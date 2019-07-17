package mongodb

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mongodb-labs/terraform-provider-mongodb/mongodb/types"
)

// WithHostSchema appends host schema to the specified schema map
func WithHostSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"host": {
			Type:     schema.TypeList,
			Required: true,
			Elem:     types.RemoteConnectionSchema,
		},
	}
}

// WithMongoDSchema appends MongoD process schema to the specified schema map
func WithMongoDSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"mongod": {
			Type:     schema.TypeList,
			Required: true,
			Elem:     types.ProcessConfigSchema,
		},
	}
}

// WithOpsManagerSchema appends OpsManager schema to the specified schema map
func WithOpsManagerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"opsmanager": {
			Type:     schema.TypeList,
			Required: true,
			Elem:     types.OpsManagerConfigSchema,
		},
	}
}

// WithAutomationSchema appends AutomationAgentConfigSchema schema to the specified schema map
func WithAutomationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"automation": {
			Type:     schema.TypeList,
			Required: true,
			Elem:     types.AutomationAgentConfigSchema,
		},
	}
}
