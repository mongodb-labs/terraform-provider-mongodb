package types

import (
	"path"

	"github.com/hashicorp/terraform/helper/schema"
)

// ProcessConfig holder for mongodb process parameters
type ProcessConfig struct {
	Binary                string  `json:"binary,omitempty"`
	WorkDir               string  `json:"workdir,omitempty"`
	Port                  int     `json:"port,string,omitempty"`
	BindIP                string  `json:"bindip,omitempty"`
	DbPath                string  `json:"dbpath,omitempty"`
	WiredTigerCacheSizeGB float64 `json:"wt_cachesize_gb,string,omitempty"`
	LogPath               string  `json:"logpath,omitempty"`
}

// ReadProcessConfig parses a singleton list of ProcessConfigSchema resources as a ProcessConfig type
func ReadProcessConfig(list []interface{}) ProcessConfig {
	// read the connection params
	cfg := &ProcessConfig{}
	data := list[0].(map[string]interface{})
	if v, ok := ReadString(data, "binary"); ok {
		cfg.Binary = v
	}
	if v, ok := ReadString(data, "workdir"); ok {
		cfg.WorkDir = v
	}
	if v, ok := ReadInt(data, "port"); ok {
		cfg.Port = v
	}
	if v, ok := ReadString(data, "bindip"); ok {
		cfg.BindIP = v
	}
	if v, ok := ReadString(data, "dbpath"); ok {
		cfg.DbPath = v
	}
	if v, ok := ReadFloat(data, "wt_cachesize_gb"); ok {
		cfg.WiredTigerCacheSizeGB = v
	}
	if v, ok := ReadString(data, "logpath"); ok {
		cfg.LogPath = v
	}
	return *cfg
}

// ProcessConfigSchema holds a minimal set of parameters required to start a MongoDB process
var ProcessConfigSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"binary": {
			Type:     schema.TypeString,
			Required: true,
		},
		"workdir": {
			Type:     schema.TypeString,
			Required: true,
		},
		"bindip": {
			Type:     schema.TypeString,
			Required: true,
		},
		"port": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  27017,
		},
		"dbpath": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "data",
		},
		"wt_cachesize_gb": {
			Type:     schema.TypeFloat,
			Optional: true,
		},
		"logpath": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "mongod.log",
		},
	},
}

// ConfigFilename returns the path to the process's config filename
func (cfg ProcessConfig) ConfigFilename() string {
	return path.Join(cfg.WorkDir, "mongod.conf")
}
