package mongodb

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
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
	sshClient, err := NewSSHClient(providerConfig, conn)
	if err != nil {
		return fmt.Errorf("could not create a SSH client: %v", err)
	}

	// create the working directory and set the appropriate permissions
	cmd := fmt.Sprintf("mkdir -p %[1]s && chown $(whoami) %[1]s && chmod 0775 %[1]s", automationConfig.WorkDir)
	ssh.PanicOnError(sshClient.RunCommand(conn.SudoPrefix(cmd)))

	// download the automation agent binary on the remote host
	// TODO(mihaibojin): Support more architectures than only 'linux_x86_64'
	filename := fmt.Sprintf("mongodb-mms-automation-agent-%s.linux_x86_64.tar.gz", automationConfig.Version)
	cmd = fmt.Sprintf("curl -O \"%s/download/agent/automation/%s\"", automationConfig.MMSBaseURL, filename)
	ssh.PanicOnError(sshClient.RunCommand(conn.SudoPrefix(cmd)))

	// unpack the binary
	cmd = fmt.Sprintf("tar -C %s -xvzf %s --strip 1", automationConfig.WorkDir, filename)
	ssh.PanicOnError(sshClient.RunCommand(cmd))
	log.Printf("[DEBUG] unpacked the binary in: %s", automationConfig.WorkDir)

	// modify automation agent config: baseUrl, ApiKey, and projectID must be set in the file along with any specified overrides
	err =
		updatePropertiesFile(sshClient, conn, automationConfig.ConfigFilename(), func(props *types.PropertiesFile) {
			props.SetPropertyValue(automationConfig.GetAutomationConfigTag("MMSGroupID"), automationConfig.MMSGroupID)
			props.SetComments(automationConfig.GetAutomationConfigTag("MMSGroupID"), []string{"", commentString, ""})
			props.SetPropertyValue(automationConfig.GetAutomationConfigTag("MMSAgentAPIKey"), automationConfig.MMSAgentAPIKey)
			props.SetPropertyValue(automationConfig.GetAutomationConfigTag("MMSBaseURL"), automationConfig.MMSBaseURL)
			props.SetPropertyValue("logFile", automationConfig.LogFilename())
		})
	util.PanicOnNonNilErr(err)

	// start the automation agent
	cmd = fmt.Sprintf("nohup %[1]s/mongodb-mms-automation-agent --config=%[2]s >> %[1]s/automation-agent-fatal.log 2>&1 & sleep 1",
		automationConfig.WorkDir, automationConfig.ConfigFilename())
	ssh.PanicOnError(sshClient.RunCommand(conn.SudoPrefix(cmd)))

	// wait for the Automation Agent to start
	if err := ssh.WaitForService(ssh.NewServiceStatusChecker(sshClient), "mongodb-mms-automation-agent"); err != nil {
		return fmt.Errorf("failed waiting for mongodb-mms-automation-agent to start: %v", err)
	}
	log.Printf("[DEBUG] confirmed the AA service is running...")

	cmd = fmt.Sprintf("(ps -ef | grep -q mongodb-mms-automation-agent && echo \"Started\" ) || echo \"Stopped\"")
	res := sshClient.RunCommand(conn.SudoPrefix(cmd))
	ssh.PanicOnError(res)
	if strings.Contains(res.Stdout, "Stopped") {
		return fmt.Errorf("error starting automation agent")
	}
	log.Printf("[DEBUG] started automation agent from configuration file %s", automationConfig.ConfigFilename())

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
