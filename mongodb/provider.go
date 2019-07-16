package mongodb

import (
	"github.com/mongodb-labs/terraform-provider-mongodb/mongodb/types"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider for MongoDB resources
func Provider() terraform.ResourceProvider {
	providerSchema, _ := getMergedProviderSchema()
	return &schema.Provider{
		Schema: providerSchema,
		ResourcesMap: map[string]*schema.Resource{
			"mongodb_process":          resourceMdbProcess(),
			"mongodb_opsmanager":       resourceMdbOpsManager(),
			"mongodb_automation_agent": resourceAutomationAgent(),
		},
		ConfigureFunc: providerConfigure,
	}
}

// getMergedProviderSchema defines all the components of our Provider schema
func getMergedProviderSchema() (map[string]*schema.Schema, error) {
	return types.StrictUnion(
		types.SSHBastionSchema(),
		types.SSHAgentSchema(),
	)
}

// providerConfigure configures the provider by parsing the passed resource data
func providerConfigure(data *schema.ResourceData) (interface{}, error) {
	bastion := types.ReadSSHBastionSchema(data)
	agent := types.ReadSSHAgentSchema(data)

	return ProviderConfig{
		Bastion: bastion,
		Agent:   agent,
	}, nil
}
