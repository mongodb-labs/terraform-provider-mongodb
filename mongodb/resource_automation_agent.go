package mongodb

import (
	"fmt"
	"log"
	"path"
	"path/filepath"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mongodb-labs/terraform-provider-mongodb/mongodb/opsmanagerconfig"
	"github.com/mongodb-labs/terraform-provider-mongodb/mongodb/ssh"
	"github.com/mongodb-labs/terraform-provider-mongodb/mongodb/types"
	"github.com/mongodb-labs/terraform-provider-mongodb/mongodb/util"
)

func resourceAutomationAgent() *schema.Resource {
	resourceSchema := types.NewSchemaMap(WithHostSchema, WithAutomationSchema)

	return &schema.Resource{
		Create: resourceMdbAutomationAgentCreate,
		Read:   resourceMdbAutomationAgentRead,
		Update: resourceMdbAutomationAgentUpdate,
		Delete: resourceMdbAutomationAgentDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(util.LongCreationTimeout),
			Read:   schema.DefaultTimeout(util.DefaultTimeout),
			Update: schema.DefaultTimeout(util.LongCreationTimeout),
			Delete: schema.DefaultTimeout(util.DefaultTimeout),
		},
		Schema: resourceSchema,
		Importer: &schema.ResourceImporter{
			State: func(data *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				return []*schema.ResourceData{data}, nil
			},
		},
	}
}

func resourceMdbAutomationAgentCreate(data *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(ProviderConfig)

	// read host params
	host := data.Get("host").([]interface{})
	conn := types.ReadRemoteConnection(host)
	data.SetId(conn.ToJSON())

	// read process config
	automation := data.Get("automation").([]interface{})
	automationConfig := types.ReadAutomationAgentConfig(automation)

	// create a SSH connection to the remote host
	client, err := NewSSHClient(providerConfig, conn)
	if err != nil {
		return fmt.Errorf("could not create a SSH client: %v", err)
	}

	// attempt to create directories if not already present
	cmd := fmt.Sprintf("mkdir -p %s %s", automationConfig.AgentDir, automationConfig.LogPath)
	ssh.PanicOnError(client.RunCommand(conn.SudoPrefix(cmd)))

	// Set correct permissions for directories
	cmd = fmt.Sprintf("chown `whoami` %s %s", automationConfig.AgentDir, automationConfig.LogPath)
	ssh.PanicOnError(client.RunCommand(conn.SudoPrefix(cmd)))

	// download the automation agent binary on the local host
	localFile, err := util.DownloadFile(automationConfig.Binary)
	util.PanicOnNonNilErr(err)
	defer util.LogError(localFile.Close)
	log.Printf("[DEBUG] downloaded automation agent binary to: %s", localFile.Name())

	// upload the binary
	remoteFilePath := path.Join(automationConfig.AgentDir, filepath.Base(localFile.Name()))
	ssh.PanicOnError(client.UploadFile(remoteFilePath, localFile))

	// unpack the binary
	cmd = fmt.Sprintf("tar -C %s -xvzf %s --strip 1", automationConfig.AgentDir, remoteFilePath)
	ssh.PanicOnError(client.RunCommand(cmd))
	log.Printf("[DEBUG] unpacked the binary in: %s", automationConfig.AgentDir)

	// TODO: get the API key and projectID from client
	// modify automation agent config: baseUrl, ApiKey, and projectID must be set in the file along with any specified overrides
	err =
		updatePropertiesFile(client, conn, automationConfig.ConfigFilename(), func(props *opsmanagerconfig.PropertiesFile) {
			props.SetPropertyValue(automationConfig.GetAutomationConfigTag("ProjectID"), automationConfig.ProjectID)
			props.SetComments(automationConfig.GetAutomationConfigTag("ProjectID"), []string{"", commentString, ""})

			props.SetPropertyValue(automationConfig.GetAutomationConfigTag("APIKey"), automationConfig.APIKey)
			props.SetComments(automationConfig.GetAutomationConfigTag("APIKey"), []string{"", commentString})

			props.SetPropertyValue(automationConfig.GetAutomationConfigTag("BaseURL"), automationConfig.BaseURL)
			props.SetComments(automationConfig.GetAutomationConfigTag("BaseURL"), []string{"", commentString})
			for prop, val := range automationConfig.Overrides {
				props.SetPropertyValue(prop, val.(string))
			}
		})
	util.PanicOnNonNilErr(err)

	// start the automation agent
	cmd = fmt.Sprintf("nohup ./mongodb-mms-automation-agent --config=%s >> %s/automation-agent-fatal.log 2>&1 &", automationConfig.ConfigFilename(), automationConfig.LogPath)
	ssh.PanicOnError(client.RunCommand(cmd))
	log.Printf("[DEBUG] started automation agent from configuration file %s", automationConfig.ConfigFilename())

	// TODO check if it's running

	return resourceMdbAutomationAgentRead(data, meta)
}

func resourceMdbAutomationAgentRead(data *schema.ResourceData, meta interface{}) error {
	// TODO implement
	return nil
}

// If the Update callback returns with or without an error, the full state is saved.
// If the ID becomes blank, the resource is destroyed (even within an update, though this shouldn't happen except in error scenarios).
// Partial mode is a mode that can be enabled by a callback that tells Terraform that it is possible for partial state to occur.
// When this mode is enabled, the provider must explicitly tell Terraform what is safe to persist and what is not.
func resourceMdbAutomationAgentUpdate(data *schema.ResourceData, meta interface{}) error {
	return resourceMdbAutomationAgentRead(data, meta)
}

// If the Destroy callback returns without an error, the resource is assumed to be destroyed, and all state is removed.
// If the Destroy callback returns with an error, the resource is assumed to still exist, and all prior state is preserved.
// If the resource is already destroyed, this should not return an error.
// This allows Terraform users to manually delete resources without breaking Terraform.
func resourceMdbAutomationAgentDelete(data *schema.ResourceData, meta interface{}) error {
	return nil
}
