package types

import (
	"path"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

// OpsManagerConfig holder for Ops Manager config
type OpsManagerConfig struct {
	Binary            string                 `json:"binary,omitempty"`
	WorkDir           string                 `json:"workdir,omitempty"`
	MongoURI          string                 `json:"mongo_uri,omitempty" opsmanager:"mongo.mongoUri"`
	EncryptionKey     string                 `json:"encryption_key,omitempty"` // /etc/mongodb-mms/gen.key
	Port              int                    `json:"port,string,omitempty" opsmanager:"BASE_PORT"`
	CentralURL        string                 `json:"central_url,omitempty" opsmanager:"mms.centralUrl"`
	Overrides         map[string]interface{} `json:"overrides,omitempty"`
	RegisterFirstUser bool                   `json:"register_first_user,omitempty"`
	FirstUserPassword string                 `json:"first_user_password,omitempty"`
	OpsManagerPort    int                    `json:"ops_manager_port,omitempty"`
	MMSGroupID        string                 `json:"mms_group_id,omitempty" automation:"mmsGroupId"`
	MMSAgentAPIKey    string                 `json:"mms_agent_api_key,omitempty" automation:"mmsApiKey"`
	Username          string                 `json:"om_user_username,omitempty"`
	Firstname         string                 `json:"om_user_firstname,omitempty"`
	Lastname          string                 `json:"om_user_lastname,omitempty"`
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
	if v, ok := ReadBool(data, "register_first_user"); ok {
		cfg.RegisterFirstUser = v
	}
	if v, ok := ReadString(data, "first_user_password"); ok {
		cfg.FirstUserPassword = v
	}
	if v, ok := ReadInt(data, "ops_manager_port"); ok {
		cfg.OpsManagerPort = v
	}
	if v, ok := ReadString(data, "mms_group_id"); ok {
		cfg.MMSGroupID = v
	}
	if v, ok := ReadString(data, "mms_agent_api_key"); ok {
		cfg.MMSAgentAPIKey = v
	}
	if v, ok := ReadString(data, "om_user_username"); ok {
		cfg.Username = v
	}
	if v, ok := ReadString(data, "om_user_firstname"); ok {
		cfg.Firstname = v
	}
	if v, ok := ReadString(data, "om_user_lastname"); ok {
		cfg.Lastname = v
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
			// TODO(mihaibojin): validate: ensure central_url and mongo_uri cannot be set here
		},
		"register_first_user": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"first_user_password": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"ops_manager_port": {
			Type:     schema.TypeInt,
			Optional: true,
		}, "ops_user_username": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "admin",
		}, "ops_user_firstname": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "adminFirst",
		}, "ops_user_lastname": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "adminLast",
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
