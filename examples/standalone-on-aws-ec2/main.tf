# Configure providers
provider "aws" {
  region                  = "eu-west-1"
  shared_credentials_file = "~/.aws/credentials"
}


# Declare variables and outputs
variable "key_name" {
  type    = "string"
  default = "terraform_ssh_key"
}

variable "aws_ssh_username" {
  type    = "string"
  default = "ec2-user"
}

variable "security_group_id" {
  type = "string"
}

variable "subnet_id" {
  type = "string"
}
output "ssh_private_key" {
  # Export with $(terraform output ssh_private_key)
  value       = tls_private_key.ssh_credentials.private_key_pem
  description = "The ssh private key used to connect to the instance."
  sensitive   = true
}
output "hostname-mdb0-0" {
  value       = aws_instance.mdb0-0.public_dns
  description = "The AWS hostname's public DNS name"
}
output "global_owner_password" {
  # Export with $(terraform output global_owner_password)
  value       = random_string.globalownerpassword.result
  description = "The password used for the auto-generated global owner account."
  sensitive   = true
}

# Deploy an AWS EC2 AMI

## Data sources
data "aws_ami" "base_ami" {
  most_recent = true

  filter {
    name = "name"
    values = [
      "RHEL-7.6_HVM_GA*"
    ]
  }

  filter {
    name = "virtualization-type"
    values = [
      "hvm"
    ]
  }

  owners = [
    # RedHat
    "309956199498"
  ]
}

## Resources
resource "tls_private_key" "ssh_credentials" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "public_ssh_key" {
  key_name   = var.key_name
  public_key = tls_private_key.ssh_credentials.public_key_openssh
}

resource "aws_instance" "mdb0-0" {
  ami           = data.aws_ami.base_ami.id
  instance_type = "t2.large"
  key_name      = aws_key_pair.public_ssh_key.key_name
  vpc_security_group_ids = [
  var.security_group_id]
  subnet_id                   = var.subnet_id
  associate_public_ip_address = true

  tags = {
    Name = "Ops Manager via Terraform Provider"
    dnr  = true
  }

  provisioner "remote-exec" {
    connection {
      type        = "ssh"
      host        = self.public_ip
      user        = var.aws_ssh_username
      private_key = tls_private_key.ssh_credentials.private_key_pem
    }

    inline = [
      # ensure the instance is actually ready and log the time
      "echo ready_at=$(date -u +'%Y-%m-%dT%H:%M:%S.%3N%z') >> instance.log",
      # the OM 4.0 AA looks for libsasl2.so.2; link it
      "sudo ln -s /usr/lib64/libsasl2.so.3 /usr/lib64/libsasl2.so.2 || echo 'Could not link libsasl2.so'"
    ]
  }
}


# Deploy a MongoDB standalone
locals {
  appdb_bind_ip = "127.0.0.1"
  ssh_port      = 22
}
resource "mongodb_process" "mdb_standalone" {
  host {
    user        = var.aws_ssh_username
    hostname    = aws_instance.mdb0-0.public_ip
    port        = local.ssh_port
    private_key = tls_private_key.ssh_credentials.private_key_pem
  }

  mongod {
    binary  = "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-rhel70-4.0.10.tgz"
    bindip  = local.appdb_bind_ip
    port    = 27017
    workdir = "/opt/mongodb"
  }
}


# Deploy a single instance of Ops Manager
resource "random_string" "encryptionkey" {
  length  = 24
  special = true
}
resource "random_string" "globalownerpassword" {
  length      = 12
  min_lower   = 1
  min_numeric = 1
  min_special = 1
}
locals {
  ops_manager_port = 9080
}
resource "mongodb_opsmanager" "opsman" {
  host {
    user        = var.aws_ssh_username
    hostname    = aws_instance.mdb0-0.public_ip
    port        = local.ssh_port
    private_key = tls_private_key.ssh_credentials.private_key_pem
  }

  opsmanager {
    binary                = "https://downloads.mongodb.com/on-prem-mms/rpm/mongodb-mms-4.0.13.50537.20190703T1029Z-1.x86_64.rpm"
    workdir               = "/opt/mongodb"
    mongo_uri             = "mongodb://${local.appdb_bind_ip}:${mongodb_process.mdb_standalone.mongod.0.port}/"
    encryption_key        = random_string.encryptionkey.result
    port                  = local.ops_manager_port
    external_port         = local.ops_manager_port
    central_url           = "http://${mongodb_process.mdb_standalone.host.0.hostname}:${local.ops_manager_port}"
    register_global_owner = true
    global_owner_username = "admin"
    global_owner_password = random_string.globalownerpassword.result

    overrides = {
      "mms.ignoreInitialUiSetup"      = "true"
      "mms.fromEmailAddr"             = "noreply@example.com"
      "mms.replyToEmailAddr"          = "noreply@example.com"
      "mms.adminEmailAddr"            = "noreply@example.com"
      "mms.mail.transport"            = "smtp"
      "mms.mail.hostname"             = "localhost"
      "mms.mail.port"                 = "25"
      "mms.mail.ssl"                  = "false"
      "automation.versions.directory" = "/data/automation/mongodb-releases"
      "automation.versions.source"    = "mongodb"
    }
  }
}


# Deploy an Automation Agent on the Ops Manager host
resource "mongodb_automation_agent" "automation_agent" {
  host {
    user        = var.aws_ssh_username
    hostname    = aws_instance.mdb0-0.public_ip
    port        = local.ssh_port
    private_key = tls_private_key.ssh_credentials.private_key_pem
  }

  automation {
    mms_base_url      = mongodb_opsmanager.opsman.opsmanager[0].central_url
    mms_group_id      = mongodb_opsmanager.opsman.opsmanager[0].mms_group_id
    mms_agent_api_key = mongodb_opsmanager.opsman.opsmanager[0].mms_agent_api_key
  }
}
