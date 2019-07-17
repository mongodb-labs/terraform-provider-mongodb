package types

import (
	"path"
	"reflect"

	"github.com/hashicorp/terraform/helper/schema"
)

// AutomationAgentConfig holder for Automation Agent Config
type AutomationAgentConfig struct {
	BaseURL   string                 `json:"baseurl,omitempty" automation:"mmsBaseUrl"`
	WorkDir   string                 `json:"workdir,omitempty"`
	Version   string                 `json:"version,omitempty"`
	LogPath   string                 `json:"logpath,omitempty"`
	ProjectID string                 `json:"project_id,omitempty" automation:"mmsGroupId"`
	APIKey    string                 `json:"agent_api_key,omitempty" automation:"mmsApiKey"`
	Overrides map[string]interface{} `json:"overrides,omitempty"`
}

// ReadAutomationAgentConfig parses a singleton list of AutomationAgentConfigSchema resources as a AutomationAgentConfig type
func ReadAutomationAgentConfig(list []interface{}) AutomationAgentConfig {
	// read the connection params
	cfg := &AutomationAgentConfig{}
	data := list[0].(map[string]interface{})
	if v, ok := ReadString(data, "agentdir"); ok {
		cfg.WorkDir = v
	}
	if v, ok := ReadString(data, "logpath"); ok {
		cfg.LogPath = v
	}
	if v, ok := ReadString(data, "baseurl"); ok {
		cfg.BaseURL = v
	}
	if v, ok := ReadString(data, "project_id"); ok {
		cfg.ProjectID = v
	}
	if v, ok := ReadString(data, "version"); ok {
		cfg.Version = v
	}
	if v, ok := ReadString(data, "agent_api_key"); ok {
		cfg.APIKey = v
	}
	if v, ok := ReadStringMap(data, "overrides"); ok {
		cfg.Overrides = v
	}
	return *cfg
}

// AutomationAgentConfigSchema holds a minimal set of parameters required to start an automation agent
var AutomationAgentConfigSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"baseurl": {
			Type:     schema.TypeString,
			Required: true,
		},
		"version": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"project_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"agent_api_key": {
			Type:     schema.TypeString,
			Required: true,
		},
		"logpath": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "/var/log/mongodb-mms-automation",
		},
		"workdir": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "/var/lib/mongodb-mms-automation",
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
