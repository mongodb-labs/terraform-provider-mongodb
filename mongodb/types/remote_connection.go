package types

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

// RemoteConnection holder for remote connection parameters
type RemoteConnection struct {
	User        string `json:"user,omitempty"`
	Hostname    string `json:"hostname,omitempty"`
	Port        int    `json:"port,string,omitempty"`
	PreventSudo bool   `json:"prevent_sudo,string,omitempty"`
	PrivateKey  string `json:"-"`
	HostKey     string `json:"-"`
}

// ReadRemoteConnectionFromString unmarshalls a JSON string into a RemoteConnection struct
func ReadRemoteConnectionFromString(data string) (*RemoteConnection, error) {
	result := &RemoteConnection{}
	if err := json.Unmarshal([]byte(data), result); err != nil {
		return nil, err
	}

	return result, nil
}

// ToJSON marshalls the struct to a JSON string
func (r RemoteConnection) ToJSON() string {
	data, err := json.Marshal(r)
	if err != nil {
		log.Fatalf("Could not marshall object to JSON: %v", r)
	}

	return string(data)
}

// SudoPrefix prefixes the specified command with 'sudo', if the RemoteConnection allows it and the user is not root
func (r RemoteConnection) SudoPrefix(cmd string) string {
	if !r.PreventSudo && r.User != "root" {
		return fmt.Sprintf("sudo %s", cmd)
	}

	return cmd
}

// ReadRemoteConnection parses a singleton list of RemoteConnectionSchema resources as a RemoteConnection type
func ReadRemoteConnection(list []interface{}) RemoteConnection {
	// read the connection params
	conn := &RemoteConnection{}
	data := list[0].(map[string]interface{})
	if v, ok := ReadString(data, "user"); ok {
		conn.User = v
	}
	if v, ok := ReadString(data, "hostname"); ok {
		conn.Hostname = v
	}
	if v, ok := ReadInt(data, "port"); ok {
		conn.Port = v
	}
	if v, ok := ReadBool(data, "prevent_sudo"); ok {
		conn.PreventSudo = v
	}
	if v, ok := ReadString(data, "private_key"); ok {
		conn.PrivateKey = v
	}
	if v, ok := ReadString(data, "host_key"); ok {
		conn.HostKey = v
	}
	return *conn
}

// RemoteConnectionSchema holds parameters used to initialize a remote SSH connection
var RemoteConnectionSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"user": {
			Type:     schema.TypeString,
			Required: true,
		},
		"hostname": {
			Type:     schema.TypeString,
			Required: true,
		},
		"port": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"prevent_sudo": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"private_key": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		"host_key": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
	},
}
