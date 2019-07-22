TEST_NAMESPACE=mongodb-terraform
TEST_CONTAINER=qa_ubuntu1804

default: apply

clean:
	@echo "==> Removing any existing terraform configs and state files..."
	@rm -rf .terraform terraform.tfstate

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

apply: init plan
	@echo "==> Applying the terraform plan..."

.PHONY: init plan apply
