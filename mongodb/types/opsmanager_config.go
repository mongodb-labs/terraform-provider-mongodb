package types

import (
	"path"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

// OpsManagerConfig holder for Ops Manager config
type OpsManagerConfig struct {
	Binary              string                 `json:"binary,omitempty"`
	WorkDir             string                 `json:"workdir,omitempty"`
	MongoURI            string                 `json:"mongo_uri,omitempty" opsmanager:"mongo.mongoUri"`
	EncryptionKey       string                 `json:"encryption_key,omitempty"` // /etc/mongodb-mms/gen.key
	Port                int                    `json:"port,omitempty" opsmanager:"BASE_PORT"`
	CentralURL          string                 `json:"central_url,omitempty" opsmanager:"mms.centralUrl"`
	Overrides           map[string]interface{} `json:"overrides,omitempty"`
	RegisterGlobalOwner bool                   `json:"register_global_owner,omitempty"`
	GlobalOwnerUsername string                 `json:"global_owner_username,omitempty"`
	GlobalOwnerPassword string                 `json:"global_owner_password,omitempty"`
	ExternalPort        int                    `json:"external_port,omitempty"`
	MMSGroupID          string                 `json:"mms_group_id,omitempty" automation:"mmsGroupId"`
	MMSAgentAPIKey      string                 `json:"mms_agent_api_key,omitempty" automation:"mmsApiKey"`
}

// ReadOpsManagerConfig parses a singleton list of OpsManagerConfigSchema resources as a OpsManagerConfig type
func ReadOpsManagerConfig(list []interface{}) OpsManagerConfig {
	// read the connection params
	cfg := &OpsManagerConfig{}
	data := list[0].(map[string]interface{})
	if v, ok := ReadString(data, "binary"); ok {
		cfg.Binary = v
	}
	if v, ok := ReadString(data, "workdir"); ok {
		cfg.WorkDir = v
	}
	if v, ok := ReadString(data, "mongo_uri"); ok {
		cfg.MongoURI = v
	}
	if v, ok := ReadString(data, "encryption_key"); ok {
		cfg.EncryptionKey = v
	}
	if v, ok := ReadInt(data, "port"); ok {
		cfg.Port = v
	}
	if v, ok := ReadString(data, "central_url"); ok {
		cfg.CentralURL = v
	}
	if v, ok := ReadStringMap(data, "overrides"); ok {
		cfg.Overrides = v
	}
	if v, ok := ReadBool(data, "register_global_owner"); ok {
		cfg.RegisterGlobalOwner = v
	}
	if v, ok := ReadString(data, "global_owner_username"); ok {
		cfg.GlobalOwnerUsername = v
	}
	if v, ok := ReadString(data, "global_owner_password"); ok {
		cfg.GlobalOwnerPassword = v
	}
	if v, ok := ReadInt(data, "external_port"); ok {
		cfg.ExternalPort = v
	}
	if v, ok := ReadString(data, "mms_group_id"); ok {
		cfg.MMSGroupID = v
	}
	if v, ok := ReadString(data, "mms_agent_api_key"); ok {
		cfg.MMSAgentAPIKey = v
	}
	return *cfg
}

// OpsManagerConfigSchema holds a minimal set of parameters required to start a MongoDB process
var OpsManagerConfigSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"binary": {
			Type:     schema.TypeString,
			Required: true,
		},
		"workdir": {
			Type:     schema.TypeString,
			Required: true,
		},
		"mongo_uri": {
			Type:     schema.TypeString,
			Required: true,
		},
		"encryption_key": {
			Type:     schema.TypeString,
			Required: true,
		},
		"port": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  8080,
		},
		"central_url": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"overrides": {
			Type:     schema.TypeMap,
			Optional: true,
			// TODO(mihaibojin): validate: central_url and mongo_uri should not be set here
		},
		"register_global_owner": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"global_owner_username": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "admin",
		},
		"global_owner_password": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"external_port": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"mms_group_id": {
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
		},
		"mms_agent_api_key": {
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
		},
	},
}

// ConfigOverrideFilename returns the path to the config override filename
func (cfg OpsManagerConfig) ConfigOverrideFilename() string {
	// if Ops Manager was installed from a tar.gz, use the working directory
	baseDir := "/opt/mongodb"
	if strings.HasSuffix(cfg.Binary, "tar.gz") || strings.HasSuffix(cfg.Binary, "tgz") {
		baseDir = cfg.WorkDir
	}

	return path.Join(baseDir, "mms", "conf", "conf-mms.properties")
}

// SysConfigFilename returns the path to the sysconfig filename
func (cfg OpsManagerConfig) SysConfigFilename() string {
	// if Ops Manager was installed from a tar.gz, use the working directory
	baseDir := "/opt/mongodb"
	if strings.HasSuffix(cfg.Binary, "tar.gz") || strings.HasSuffix(cfg.Binary, "tgz") {
		baseDir = cfg.WorkDir
	}

	return path.Join(baseDir, "mms", "conf", "mms.conf")
}

// GetOpsManagerTag returns the specified opsmanager tag
func (cfg OpsManagerConfig) GetOpsManagerTag(fieldName string) string {
	t := reflect.TypeOf(cfg)
	field, _ := t.FieldByName(fieldName)
	return field.Tag.Get("opsmanager")
}
