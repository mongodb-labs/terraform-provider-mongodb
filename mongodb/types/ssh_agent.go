package types

import (
	"github.com/MihaiBojin/terraform-provider-mongodb/mongodb/ssh"
	"github.com/hashicorp/terraform/helper/schema"
)

// SSHAgentSchema constructs a terraform schema map representing SSH agent configuration params
func SSHAgentSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"agent": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "",
		},
		"agent_identity": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "",
		},
	}
}

// ReadSSHAgentSchema reads SSH agent configuration from the passed schema.ResourceData struct
func ReadSSHAgentSchema(data *schema.ResourceData) ssh.Agent {
	agent := data.Get("agent").(bool)
	agentIdentity := data.Get("agent_identity").(string)

	return ssh.Agent{
		Agent:         agent,
		AgentIdentity: agentIdentity,
	}
}
