package types

import (
	"github.com/hashicorp/terraform/helper/schema"
	"path"
	"reflect"
)

type AutomationAgentConfig struct {
	Binary    string                 `json:"binary,omitempty"`
	BaseUrl   string                 `json:"baseurl,omitempty" automation:"mmsBaseUrl"`
	AgentDir  string                 `json:"agentdir,omitempty"`
	LogPath   string                 `json:"logpath,omitempty"`
	GroupId   string                 `json:"group_id,omitempty" automation:"mmsGroupId"`
	ApiKey    string                 `json:"api_key,omitempty" automation:"mmsApiKey"`
	Overrides map[string]interface{} `json:"overrides,omitempty"`
}

// ReadAutomationAgentConfig parses a singleton list of AutomationAgentConfigSchema resources as a AutomationAgentConfig type
func ReadAutomationAgentConfig(list []interface{}) AutomationAgentConfig {
	// read the connection params
	cfg := &AutomationAgentConfig{}
	data := list[0].(map[string]interface{})
	if v, ok := ReadString(data, "binary"); ok {
		cfg.Binary = v
	}
	if v, ok := ReadString(data, "agentdir"); ok {
		cfg.AgentDir = v
	}
	if v, ok := ReadString(data, "logpath"); ok {
		cfg.LogPath = v
	}
	if v, ok := ReadString(data, "baseurl"); ok {
		cfg.BaseUrl = v
	}
	if v, ok := ReadString(data, "group_id"); ok {
		cfg.GroupId = v
	}
	if v, ok := ReadString(data, "api_key"); ok {
		cfg.ApiKey = v
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
		"group_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"api_key": {
			Type:     schema.TypeString,
			Required: true,
		},
		"binary": {
			Type:     schema.TypeString,
			Required: true,
		},
		"logpath": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "/var/log/mongodb-mms-automation",
		},
		"agentdir": {
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
	return path.Join(cfg.AgentDir, "local.config")
}

// GetAutomationConfigTag given a valid AutomationConfig struct field name, returns the specified automation tag
func (cfg AutomationAgentConfig) GetAutomationConfigTag(fieldName string) string {
	t := reflect.TypeOf(cfg)
	field, _ := t.FieldByName(fieldName)
	return field.Tag.Get("automation")
}
