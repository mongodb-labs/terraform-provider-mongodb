package ssh

import (
	"encoding/json"
	"github.com/MihaiBojin/terraform-provider-mongodb/mongodb/util"

	"github.com/hashicorp/terraform/terraform"
)

// Connection holds all connection parameters used to connect to a host via SSH
type Connection struct {
	User       string `json:"user,omitempty"`
	Password   string `json:"password,omitempty"`
	PrivateKey string `json:"private_key,omitempty"`
	Host       string `json:"host,omitempty"`
	HostKey    string `json:"host_key,omitempty"`
	Port       int    `json:"port,string,omitempty"`
	Bastion
	Agent
}

// WithHostParams accepts the minimal set of parameters required to make a connection
func WithHostParams(user string, host string, port int) func(params *Connection) error {
	return func(params *Connection) error {
		params.User = user
		params.Host = host
		params.Port = port
		return nil
	}
}

// WithPrivateKey specifies the private keys to use for connecting to the target host
func WithPrivateKey(privateKey string) func(params *Connection) error {
	return func(params *Connection) error {
		params.PrivateKey = privateKey
		return nil
	}
}

// WithHostKey specifies the host keys of the target host
func WithHostKey(hostKey string) func(params *Connection) error {
	return func(params *Connection) error {
		params.HostKey = hostKey
		return nil
	}
}

// ToMap converts *Connection to a map[string]string by marshaling the struct to json bytes and then unmarshalling into the desired output format
func (connection *Connection) toMap() (map[string]string, error) {
	var result map[string]string
	inrec, err := json.Marshal(connection)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(inrec, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// toEphemeralState constructs a terraform.EphemeralState object to be used by the SSH communicator
func (connection *Connection) toEphemeralState() *terraform.EphemeralState {
	connMap, err := connection.toMap()
	util.PanicOnNonNilErr(err)
	return &terraform.EphemeralState{ConnInfo: connMap}
}
