package types

import (
	"github.com/MihaiBojin/terraform-provider-mongodb/mongodb/ssh"
	"github.com/hashicorp/terraform/helper/schema"
)

// SSHBastionSchema constructs a terraform schema map representing SSH bastion host connection params
func SSHBastionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"bastion_user": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "",
		},
		"bastion_password": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "",
		},
		"bastion_private_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "",
		},
		"bastion_host": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "",
		},
		"bastion_host_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "",
		},
		"bastion_port": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     22,
			Description: "",
		},
	}
}

// ReadSSHBastionSchema reads bastion host configuration from the passed schema.ResourceData struct
func ReadSSHBastionSchema(data *schema.ResourceData) ssh.Bastion {
	bastionUser := data.Get("bastion_user").(string)
	bastionPassword := data.Get("bastion_password").(string)
	bastionPrivateKey := data.Get("bastion_private_key").(string)
	bastionHost := data.Get("bastion_host").(string)
	bastionHostKey := data.Get("bastion_host_key").(string)
	bastionPort := data.Get("bastion_port").(int)

	return ssh.Bastion{
		BastionUser:       bastionUser,
		BastionPassword:   bastionPassword,
		BastionPrivateKey: bastionPrivateKey,
		BastionHost:       bastionHost,
		BastionHostKey:    bastionHostKey,
		BastionPort:       bastionPort,
	}
}
