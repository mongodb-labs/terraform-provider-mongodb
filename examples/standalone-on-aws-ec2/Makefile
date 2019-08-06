default: apply

clean:
	@echo "==> Removing any existing terraform configs and state files..."
	@rm -rf .terraform terraform.tfstate terraform.tfstate.backup

init:
	@echo "==> Initializing terraform..."
	@terraform init

ifndef TF_VAR_security_group_id
	$(warning TF_VAR_security_group_id is not set in the environment. Terraform will ask you to input a value...)
endif
ifndef TF_VAR_subnet_id
	$(warning TF_VAR_subnet_id is not set in the environment. Terraform will ask you to input a value...)
endif
ifeq ($(shell test -f terraform.tfstate && echo -n yes),yes)
	$(error Pre-existing terraform state detected.  You can remove it and start fresh with "make clean"...)
endif
plan:
	@echo "==> Creating the terraform plan..."
	@terraform plan

private-key:
	@$(eval TMP_FILE = $(shell mktemp))
	@TF_LOG=ERROR terraform output ssh_private_key > $(TMP_FILE)
	@chmod 0500 $(TMP_FILE)
	@echo $(TMP_FILE)

ssh:
	@$(eval SSHHOST = $(shell TF_LOG=ERROR terraform output hostname-mdb0-1))
	@echo "==> Connecting to $(SSHHOST) ..."
	@$(eval PKEY = $(shell make private-key))
	@ssh -i $(PKEY) -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ec2-user@$(SSHHOST)

apply: init plan
	@echo "==> Applying the terraform plan..."
	@terraform apply

.PHONY: clean init plan apply private-key
