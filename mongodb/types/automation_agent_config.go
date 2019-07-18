package types

import (
	"path"
	"reflect"

	"github.com/hashicorp/terraform/helper/schema"
)

// AutomationAgentConfig holder for Automation Agent Config
type AutomationAgentConfig struct {
	MMSBaseURL     string                 `json:"mms_base_url,omitempty" automation:"mmsBaseUrl"`
	WorkDir        string                 `json:"workdir,omitempty"`
	Version        string                 `json:"version,omitempty"`
	LogPath        string                 `json:"logpath,omitempty"`
	MMSGroupID     string                 `json:"mms_group_id,omitempty" automation:"mmsGroupId"`
	MMSAgentAPIKey string                 `json:"mms_agent_api_key,omitempty" automation:"mmsApiKey"`
	Overrides      map[string]interface{} `json:"overrides,omitempty"`
}

// ReadAutomationAgentConfig parses a singleton list of AutomationAgentConfigSchema resources as a AutomationAgentConfig type
func ReadAutomationAgentConfig(list []interface{}) AutomationAgentConfig {
	// read the connection params
	cfg := &AutomationAgentConfig{}
	data := list[0].(map[string]interface{})
	if v, ok := ReadString(data, "mms_base_url"); ok {
		cfg.MMSBaseURL = v
	}
	if v, ok := ReadString(data, "mms_group_id"); ok {
		cfg.MMSGroupID = v
	}
	if v, ok := ReadString(data, "mms_agent_api_key"); ok {
		cfg.MMSAgentAPIKey = v
	}
	if v, ok := ReadString(data, "version"); ok {
		cfg.Version = v
	}
	if v, ok := ReadString(data, "workdir"); ok {
		cfg.WorkDir = v
	}
	if v, ok := ReadString(data, "logpath"); ok {
		cfg.LogPath = v
	}
	if v, ok := ReadStringMap(data, "overrides"); ok {
		cfg.Overrides = v
	}
	return *cfg
}

// AutomationAgentConfigSchema holds a minimal set of parameters required to start an automation agent
var AutomationAgentConfigSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"mms_base_url": {
			Type:     schema.TypeString,
			Required: true,
		},
		"mms_group_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"mms_agent_api_key": {
			Type:     schema.TypeString,
			Required: true,
		},
		"version": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "latest",
		},
		"workdir": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "/var/lib/mongodb-mms-automation",
		},
		"logpath": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "/var/log/mongodb-mms-automation",
		},
		"overrides": {
			Type:     schema.TypeMap,
			Optional: true,
		},
	},
}

// ConfigFilename returns the path to the process's config filename
func (cfg AutomationAgentConfig) ConfigFilename() string {
	return path.Join(cfg.WorkDir, "local.config")
}

// GetAutomationConfigTag given a valid AutomationConfig struct field name, returns the specified automation tag
func (cfg AutomationAgentConfig) GetAutomationConfigTag(fieldName string) string {
	t := reflect.TypeOf(cfg)
	field, _ := t.FieldByName(fieldName)
	return field.Tag.Get("automation")
}
