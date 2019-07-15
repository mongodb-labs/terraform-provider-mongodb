package ssh

// Bastion holds all parameters required to make an SSH connection through a bastion host
type Bastion struct {
	BastionUser       string `json:"bastion_user,omitempty"`
	BastionPassword   string `json:"bastion_password,omitempty"`
	BastionPrivateKey string `json:"bastion_private_key,omitempty"`
	BastionHost       string `json:"bastion_host,omitempty"`
	BastionHostKey    string `json:"bastion_host_key,omitempty"`
	BastionPort       int    `json:"bastion_port,string,omitempty"`
}
