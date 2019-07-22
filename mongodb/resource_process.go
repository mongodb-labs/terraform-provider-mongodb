package mongodb

import (
	"fmt"
	"log"
	"path"
	"path/filepath"

	"github.com/mongodb-labs/terraform-provider-mongodb/mongodb/config"
	"github.com/mongodb-labs/terraform-provider-mongodb/mongodb/ssh"
	"github.com/mongodb-labs/terraform-provider-mongodb/mongodb/types"
	"github.com/mongodb-labs/terraform-provider-mongodb/mongodb/util"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceMdbProcess() *schema.Resource {
	// TODO(mihaibojin): split 'mongod' schema: move binary at the top level
	// TODO(mihaibojin): split 'mongod' schema: move all other settings under 'host'
	// TODO(mihaibojin): rename 'host' to 'instance'
	// TODO(mihaibojin): rename 'process' to 'database' (or a better name that represents an unmanaged standalone or replica set)
	resourceSchema := types.NewSchemaMap(WithHostSchema, WithMongoDSchema)

	return &schema.Resource{
		Create: resourceMdbProcessCreate,
		Read:   resourceMdbProcessRead,
		Update: resourceMdbProcessUpdate,
		Delete: resourceMdbProcessDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(util.LongCreationTimeout),
			Read:   schema.DefaultTimeout(util.DefaultTimeout),
			Update: schema.DefaultTimeout(util.LongCreationTimeout),
			Delete: schema.DefaultTimeout(util.DefaultTimeout),
		},
		Schema: resourceSchema,
	}
}

// NewSSHClient build a new SSH client using the provided parameters
func NewSSHClient(pc ProviderConfig, rc types.RemoteConnection) (*ssh.Client, error) {
	return ssh.NewClient(
		ssh.WithHostParams(rc.User, rc.Hostname, rc.Port),
		WithProviderConfig(&pc),
		ssh.WithPrivateKey(rc.PrivateKey),
		ssh.WithHostKey(rc.HostKey),
	)
}

// If the Create callback returns with or without an error without an ID set using SetId, the resource is assumed to not be created, and no state is saved.
// If the Create callback returns with or without an error and an ID has been set, the resource is assumed created and all state is saved with it.
func resourceMdbProcessCreate(data *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(ProviderConfig)

	// read host params
	host := data.Get("host").([]interface{})
	conn := types.ReadRemoteConnection(host)
	data.SetId(conn.ToJSON())

	// read process config
	process := data.Get("mongod").([]interface{})
	dbConfig := types.ReadProcessConfig(process)

	// create a SSH connection to the remote host
	client, err := NewSSHClient(providerConfig, conn)
	if err != nil {
		return fmt.Errorf("could not create a SSH client: %v", err)
	}

	// create the working directory (and other dirs) and set the appropriate permissions
	dbPath := filepath.Join(dbConfig.WorkDir, dbConfig.DbPath)
	logPath := filepath.Join(dbConfig.WorkDir, dbConfig.DbPath, dbConfig.LogPath)
	cmd := fmt.Sprintf("mkdir -p %[1]s %[2]s %[3]s && chown $(whoami) %[1]s %[2]s %[3]s && chmod 0775 %[1]s %[2]s %[3]s", dbConfig.WorkDir, dbPath, path.Base(logPath))
	ssh.PanicOnError(client.RunCommand(conn.SudoPrefix(cmd)))

	// create a MongoDB configuration file
	cfg := config.NewMongoDBConfig()
	cfg.ProcessManagement.Fork = true
	cfg.Storage.DBPath = dbPath
	cfg.Net = &config.Net{
		Port:   dbConfig.Port,
		BindIP: dbConfig.BindIP,
	}
	cfg.Storage.Engine = "wiredTiger"
	if dbConfig.WiredTigerCacheSizeGB > 0 {
		cfg.Storage.WiredTiger.EngineConfig.CacheSizeGB = dbConfig.WiredTigerCacheSizeGB
	}
	cfg.SystemLog.LogAppend = true
	cfg.SystemLog.Path = logPath
	cfg.SystemLog.Destination = "file"
	cfgFile, err := cfg.SaveToTempFile("")
	util.PanicOnNonNilErr(err)
	defer util.BurnAfterReading(cfgFile)

	// upload the config file to the remote host
	remoteConfigPath := dbConfig.ConfigFilename()
	ssh.PanicOnError(client.UploadFile(remoteConfigPath, cfgFile))

	// download the MongoDB binary on the local host
	localFile, err := util.DownloadFile(dbConfig.Binary)
	util.PanicOnNonNilErr(err)
	defer util.LogError(localFile.Close)
	log.Printf("[DEBUG] downloaded binary to: %s", localFile.Name())

	// upload the binary
	remoteFilePath := path.Join(dbConfig.WorkDir, filepath.Base(localFile.Name()))
	ssh.PanicOnError(client.UploadFile(remoteFilePath, localFile))

	// unpack the binary
	cmd = fmt.Sprintf("tar -C %s -xvzf %s --strip 1", dbConfig.WorkDir, remoteFilePath)
	ssh.PanicOnError(client.RunCommand(cmd))
	log.Printf("[DEBUG] unpacked the binary in: %s", dbConfig.WorkDir)

	// start the process
	cmd = fmt.Sprintf("%s/bin/mongod -f %s || cat %s", dbConfig.WorkDir, remoteConfigPath, logPath)
	ssh.PanicOnError(client.RunCommand(conn.SudoPrefix(cmd)))
	log.Print("[DEBUG] started MongoD...")

	// check the connection
	cmd = fmt.Sprintf("%s/bin/mongo --quiet --port %d --eval 'quit()'", dbConfig.WorkDir, dbConfig.Port)
	ssh.PanicOnError(client.RunCommand(cmd))
	log.Printf("[DEBUG] Successfully connected to MongoDB on port %d", dbConfig.Port)

	return resourceMdbProcessRead(data, meta)
}

// This callback should never modify the real resource.
// If the ID is updated to blank, this tells Terraform the resource no longer exists (maybe it was destroyed out of band).
// Just like the destroy callback, the Read function should gracefully handle this case.
// https://www.terraform.io/docs/extend/writing-custom-providers.html#implementing-a-more-complex-read
func resourceMdbProcessRead(data *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(ProviderConfig)

	// read host params
	host := data.Get("host").([]interface{})
	conn := types.ReadRemoteConnection(host)

	// TODO(mihaibojin): set the correct id

	// read process config
	process := data.Get("mongod").([]interface{})
	currentConfig := types.ReadProcessConfig(process)

	// create a SSH connection to the remote host
	client, err := NewSSHClient(providerConfig, conn)
	if err != nil {
		return fmt.Errorf("could not create a SSH client: %v", err)
	}

	// load the configuration file
	result := client.RunCommand(fmt.Sprintf("cat %s", currentConfig.ConfigFilename()))
	ssh.PanicOnError(result)

	// parse the configuration into a config.MongoDB struct
	mongoDBConfig, err := config.LoadFromString(result.Stdout)
	if err != nil {
		return err
	}

	// update the resource data
	resourceData := make(map[string]interface{})
	resourceData["binary"] = currentConfig.Binary
	resourceData["port"] = mongoDBConfig.Net.Port
	resourceData["bindip"] = mongoDBConfig.Net.BindIP
	resourceData["dbpath"] = mongoDBConfig.Storage.DBPath
	resourceData["logpath"] = mongoDBConfig.SystemLog.Path
	if err := data.Set("mongod", []map[string]interface{}{resourceData}); err != nil { // convert resourceData to an array of maps with a single element
		return err
	}

	log.Print("[DEBUG] updated the MongoDB Process resource...")
	return nil
}

// If the Update callback returns with or without an error, the full state is saved.
// If the ID becomes blank, the resource is destroyed (even within an update, though this shouldn't happen except in error scenarios).
// Partial mode is a mode that can be enabled by a callback that tells Terraform that it is possible for partial state to occur.
// When this mode is enabled, the provider must explicitly tell Terraform what is safe to persist and what is not.
func resourceMdbProcessUpdate(data *schema.ResourceData, meta interface{}) error {
	return resourceMdbProcessRead(data, meta)
}

// If the Destroy callback returns without an error, the resource is assumed to be destroyed, and all state is removed.
// If the Destroy callback returns with an error, the resource is assumed to still exist, and all prior state is preserved.
// If the resource is already destroyed, this should not return an error.
// This allows Terraform users to manually delete resources without breaking Terraform.
func resourceMdbProcessDelete(data *schema.ResourceData, meta interface{}) error {
	return nil
}
