# Configure providers
provider "aws" {
  region                  = "eu-west-1"
  shared_credentials_file = "~/.aws/credentials"
}

#
# AWS AMI with MongoDB support
#

## Variables
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

## Datasources
data "aws_ami" "base_ami" {
  most_recent = true

  filter {
    name = "name"
    values = [
    "RHEL-7.6_HVM_GA*"]
  }

  filter {
    name = "virtualization-type"
    values = [
    "hvm"]
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
  instance_type = "t2.micro"
  key_name      = aws_key_pair.public_ssh_key.key_name
  vpc_security_group_ids = [
  var.security_group_id]
  subnet_id                   = var.subnet_id
  associate_public_ip_address = true

  tags = {
    Name = "RedHat 7.6 + MongoDB support"
  }

  provisioner "remote-exec" {
    connection {
      type        = "ssh"
      user        = var.aws_ssh_username
      private_key = tls_private_key.ssh_credentials.private_key_pem
    }

    inline = [
      # ensure the instance is actually ready and log the time
      "echo ready_at=$(date -u +'%Y-%m-%dT%H:%M:%S.%3N%z') >> instance.log",
    ]
  }
}

## Outputs
output "mdb0-0-connection" {
  value       = "${var.aws_ssh_username}@${aws_instance.mdb0-0.public_ip}"
  description = "The FQDN of the provisioned instance."
}

output "ssh_private_key" {
  # Export with $(terraform output ssh_private_key)
  value       = tls_private_key.ssh_credentials.private_key_pem
  description = "The ssh private key used to connect to the instance."
  sensitive   = true
}


#
# MongoDB standalone
#

resource "mongodb_process" "standalone" {
  host {
    user     = "root"
    hostname = "127.0.0.1"
    port     = 22
  }

  mongod {
    binary  = "http://downloads.mongodb.org/linux/mongodb-linux-x86_64-ubuntu1804-4.0.10.tgz"
    bindip  = "0.0.0.0"
    workdir = "/opt/mongodb"
  }
}

# TODO(mihaibojin): Complete this example (based on the Docker equivalent)
