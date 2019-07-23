default: apply

clean:
	@echo "==> Removing any existing terraform configs and state files..."
	@rm -rf .terraform terraform.tfstate terraform.tfstate.backup

init:
	@echo "==> Initializing terraform..."
	@terraform init

plan:
	@echo "==> Creating the terraform plan..."
	@terraform plan

apply: init plan
	@echo "==> Applying the terraform plan..."
	@terraform apply

.PHONY: clean init plan apply
