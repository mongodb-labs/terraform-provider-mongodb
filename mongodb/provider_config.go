package mongodb

import (
	"github.com/MihaiBojin/terraform-provider-mongodb/mongodb/ssh"
)

// ProviderConfig provider configuration object
type ProviderConfig struct {
	ssh.Bastion
	ssh.Agent
}

// WithProviderConfig helper for passing provider configuration to the SSH client via a *ssh.Connection
func WithProviderConfig(pc *ProviderConfig) func(params *ssh.Connection) error {
	return func(params *ssh.Connection) error {
		params.BastionUser = pc.BastionUser
		params.BastionHost = pc.BastionHost
		params.BastionPort = pc.BastionPort
		params.BastionPrivateKey = pc.BastionPrivateKey
		params.BastionHostKey = pc.BastionHostKey

		params.Agent.Agent = pc.Agent.Agent
		params.AgentIdentity = pc.AgentIdentity

		return nil
	}
}
