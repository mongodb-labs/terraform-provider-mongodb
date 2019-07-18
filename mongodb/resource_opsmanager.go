package mongodb

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mongodb-labs/pcgc/pkg/httpclient"
	"github.com/mongodb-labs/pcgc/pkg/opsmanager"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mongodb-labs/terraform-provider-mongodb/mongodb/ssh"
	"github.com/mongodb-labs/terraform-provider-mongodb/mongodb/types"
	"github.com/mongodb-labs/terraform-provider-mongodb/mongodb/util"
)

const commentString = "# DO NOT CHANGE - this file was generated by the MongoDB Terraform Provider"

func resourceMdbOpsManager() *schema.Resource {
	// TODO(mihaibojin): 'opsmanager' schema: move settings at the top-level
	resourceSchema := types.NewSchemaMap(WithHostSchema, WithOpsManagerSchema)

	return &schema.Resource{
		Create: resourceMdbOpsManagerCreate,
		Read:   resourceMdbOpsManagerRead,
		Update: resourceMdbOpsManagerUpdate,
		Delete: resourceMdbOpsManagerDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(util.LongCreationTimeout),
			Read:   schema.DefaultTimeout(util.DefaultTimeout),
			Update: schema.DefaultTimeout(util.LongCreationTimeout),
			Delete: schema.DefaultTimeout(util.DefaultTimeout),
		},
		Schema: resourceSchema,
	}
}

// If the Create callback returns with or without an error without an ID set using SetId, the resource is assumed to not be created, and no state is saved.
// If the Create callback returns with or without an error and an ID has been set, the resource is assumed created and all state is saved with it.
func resourceMdbOpsManagerCreate(data *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(ProviderConfig)

	host := data.Get("host").([]interface{})
	conn := types.ReadRemoteConnection(host)
	data.SetId(conn.ToJSON())

	om := data.Get("opsmanager").([]interface{})
	omConfig := types.ReadOpsManagerConfig(om)

	// create a SSH connection to the remote host
	client, err := NewSSHClient(providerConfig, conn)
	if err != nil {
		return fmt.Errorf("could not create a SSH client: %v", err)
	}

	// attempt to create the work directory
	cmd := fmt.Sprintf("mkdir -p %s", omConfig.WorkDir)
	ssh.PanicOnError(client.RunCommand(conn.SudoPrefix(cmd)))

	// download Ops Manager
	localFile, err := util.DownloadFile(omConfig.Binary)
	util.PanicOnNonNilErr(err)
	defer util.LogError(localFile.Close)
	log.Printf("[DEBUG] downloaded binary to: %s", localFile.Name())

	// upload the binary
	fileName := filepath.Base(localFile.Name())
	remoteFilePath := path.Join(omConfig.WorkDir, fileName)
	ssh.PanicOnError(client.UploadFile(remoteFilePath, localFile))
	log.Printf("[DEBUG] uploaded the binary to: %s", remoteFilePath)

	// install Ops Manager
	filetype := filepath.Ext(localFile.Name())
	if filetype == ".tar.gz" || filetype == ".tgz" {
		// unpack the binary
		cmd = fmt.Sprintf("tar -C %s -xvzf %s --strip 1", omConfig.WorkDir, remoteFilePath)
		ssh.PanicOnError(client.RunCommand(cmd))

		// TODO(mihaibojin): support installing from tar.gz
	} else if filetype == ".deb" {
		// install the binary
		cmd := fmt.Sprintf(conn.SudoPrefix("dpkg -i --force-confnew %s"), remoteFilePath)
		ssh.PanicOnError(client.RunCommand(cmd))
	} else if filetype == ".rpm" {
		// install the binary
		cmd := fmt.Sprintf(conn.SudoPrefix("rpm -ivh %s"), remoteFilePath)
		ssh.PanicOnError(client.RunCommand(cmd))
	} else {
		return fmt.Errorf("unknown file type: %v", filetype)
	}
	log.Print("[DEBUG] unpacked the binary on the remote host")

	// configure the property overrides (conf-mms.properties)
	err =
		updatePropertiesFile(client, conn, omConfig.ConfigOverrideFilename(), func(props *types.PropertiesFile) {
			props.SetPropertyValue(omConfig.GetOpsManagerTag("MongoURI"), omConfig.MongoURI)
			props.SetComments(omConfig.GetOpsManagerTag("MongoURI"), []string{"", commentString, ""})

			props.SetPropertyValue(omConfig.GetOpsManagerTag("CentralURL"), omConfig.CentralURL)
			props.SetComments(omConfig.GetOpsManagerTag("CentralURL"), []string{"", commentString})
			for prop, val := range omConfig.Overrides {
				props.SetPropertyValue(prop, val.(string))
			}
		})
	util.PanicOnNonNilErr(err)

	// configure the port in mms.conf
	err =
		updatePropertiesFile(client, conn, omConfig.SysConfigFilename(), func(props *types.PropertiesFile) {
			props.SetPropertyValue(omConfig.GetOpsManagerTag("Port"), strconv.Itoa(omConfig.Port))
			props.SetComments(omConfig.GetOpsManagerTag("Port"), []string{commentString, ""})
		})
	util.PanicOnNonNilErr(err)

	// upload the encryption key
	remoteEncKeyPath := "/etc/mongodb-mms/gen.key"
	remoteTempFile := "~/gen.key"
	// create the encryption key's directory
	cmd = fmt.Sprintf("mkdir -p %s", filepath.Dir(remoteEncKeyPath))
	ssh.PanicOnError(client.RunCommand(conn.SudoPrefix(cmd)))
	// store the encryption key to a temp file, always ensuring no more than 24 bytes are selected
	escaped := strings.Replace(omConfig.EncryptionKey[0:24], "'", "\\'", -1)
	encKeyFile, err := util.ReadAllIntoTempFile(strings.NewReader(escaped), "encryption-key")
	util.PanicOnNonNilErr(err)
	defer util.BurnAfterReading(encKeyFile)
	// upload the file
	ssh.PanicOnError(client.UploadFile(remoteTempFile, encKeyFile))
	// move the file to its final location and set the correct perms
	cmd = fmt.Sprintf("bash -c 'mv %[1]s %[2]s; chown mongodb-mms:mongodb-mms %[2]s; chmod 0600 %[2]s'", remoteTempFile, remoteEncKeyPath)
	ssh.PanicOnError(client.RunCommand(conn.SudoPrefix(cmd)))

	// start the Ops Manager service
	ssh.PanicOnError(client.RunCommand(conn.SudoPrefix("/etc/init.d/mongodb-mms start")))
	log.Printf("[DEBUG] started Ops Manager on port: %d", omConfig.Port)

	// wait for Ops Manager to start
	if err := ssh.WaitForOpenPort(ssh.NewOpenPortCheckerFunc(client), omConfig.Port); err != nil {
		return err
	}
	log.Printf("[DEBUG] confirmed connection to the Ops Manager port: %d", omConfig.Port)

	// create first user if option was passed
	if omConfig.RegisterFirstUser {
		// create the first user via the client with noauth
		apiURL := fmt.Sprintf("http://%s:%d", conn.Hostname, omConfig.OpsManagerPort)
		resolver := httpclient.NewURLResolverWithPrefix(apiURL, opsmanager.PublicAPIPrefix)
		apiFirstUserResp, err := createFirstUser(resolver, omConfig.FirstUserPassword)
		if err != nil {
			return fmt.Errorf("Failed to create first user: %v", err)
		}
		log.Printf("[DEBUG] Created first OM user: %s", apiFirstUserResp.User.Username)

		// create the first project via the client with digestAuth, to get projectID and agentAPIKey
		createOneProjectResp, err := createFirstProject(resolver, apiFirstUserResp.User.Username, apiFirstUserResp.APIKey)
		if err != nil {
			return fmt.Errorf("Failed to create first project using Private Cloud Go Client: %v", err)
		}
		log.Printf("[DEBUG] Created first project using the Private Cloud Go Client. ProjectId, agent API key: %s , %s", createOneProjectResp.ID, createOneProjectResp.AgentAPIKey)

		omConfig.MMSAgentAPIKey = createOneProjectResp.AgentAPIKey
		omConfig.MMSGroupID = createOneProjectResp.ID
	}

	return resourceMdbOpsManagerRead(data, meta)
}

// This callback should never modify the real resource.
// If the ID is updated to blank, this tells Terraform the resource no longer exists (maybe it was destroyed out of band).
// Just like the destroy callback, the Read function should gracefully handle this case.
func resourceMdbOpsManagerRead(data *schema.ResourceData, meta interface{}) error {
	// https://www.terraform.io/docs/extend/writing-custom-providers.html#implementing-a-more-complex-read

	// TODO(mihaibojin): set the correct id

	// ssh connect
	// attempt to connect to MongoDB via port
	ok := true

	// If the resource does not exist, inform Terraform. We want to immediately
	// return here to prevent further processing.
	if !ok {
		log.Printf("[WARN] could not connect: %s", data.Id())
		return nil
	}

	log.Print("[DEBUG] updated the Ops Manager resource...")
	return nil
}

// If the Update callback returns with or without an error, the full state is saved.
// If the ID becomes blank, the resource is destroyed (even within an update, though this shouldn't happen except in error scenarios).
// Partial mode is a mode that can be enabled by a callback that tells Terraform that it is possible for partial state to occur.
// When this mode is enabled, the provider must explicitly tell Terraform what is safe to persist and what is not.
func resourceMdbOpsManagerUpdate(data *schema.ResourceData, meta interface{}) error {
	return resourceMdbProcessRead(data, meta)
}

// If the Destroy callback returns without an error, the resource is assumed to be destroyed, and all state is removed.
// If the Destroy callback returns with an error, the resource is assumed to still exist, and all prior state is preserved.
// If the resource is already destroyed, this should not return an error.
// This allows Terraform users to manually delete resources without breaking Terraform.
func resourceMdbOpsManagerDelete(data *schema.ResourceData, meta interface{}) error {
	return nil
}

// updatePropertiesFile updates a remote property file, given a set of modifications defined in updateProps
func updatePropertiesFile(client *ssh.Client, conn types.RemoteConnection, remoteFile string, updateProps func(*types.PropertiesFile)) error {
	// back up the old file
	ssh.PanicOnError(client.RunCommand(conn.SudoPrefix(fmt.Sprintf("cp %s %s.backup", remoteFile, remoteFile))))
	log.Printf("[DEBUG] backed up: %s", remoteFile)

	// download the configuration file
	result := client.RunCommand(conn.SudoPrefix(fmt.Sprintf("cat %s", remoteFile)))
	ssh.PanicOnError(result)
	log.Printf("[DEBUG] downloaded the file from: %s", remoteFile)

	// parse the configuration into a struct and apply the updates
	config := types.NewPropertiesFile(result.Stdout)
	updateProps(config)
	configData, err := config.Write()
	util.PanicOnNonNilErr(err)

	// upload the config file to the remote host
	ssh.PanicOnError(client.UploadData(remoteFile, bufio.NewReader(strings.NewReader(configData))))
	log.Printf("[DEBUG] uploaded the config file to the remote host, at: %s", remoteFile)

	return nil
}

func createFirstUser(resolver httpclient.URLResolver, firstPassword string) (resp opsmanager.CreateFirstUserResponse, err error) {
	// initialize a client without auth
	omAPIClientNoAuth := opsmanager.NewDefaultClient(resolver)

	// create the first user
	user := opsmanager.User{Username: "firstuser", Password: firstPassword, FirstName: "first", LastName: "last"}
	return omAPIClientNoAuth.CreateFirstUser(user, "0.0.0.1/0")

}

func createFirstProject(resolver httpclient.URLResolver, username string, apiKey string) (resp opsmanager.CreateOneProjectResponse, err error) {
	// initialize a client with auth
	omAPIClientDigestAuth := opsmanager.NewClientWithDigestAuth(resolver, username, apiKey)

	// create new org/project to get the GroupID
	projectName := fmt.Sprintf("TerraformProject-%d", rand.Intn(10000000))
	return omAPIClientDigestAuth.CreateOneProject(projectName, "")
}
