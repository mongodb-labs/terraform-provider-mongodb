TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=mongodb

default: build

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

debugacc: fmtcheck
	TF_ACC=1 dlv test $(TEST) -- -test.v $(TESTARGS)

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -s -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

lint:
	golint -set_exit_status ./...
	GOGC=30 golangci-lint run ./...

vendor-status:
	@govendor status

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build install test testacc vet fmt fmtcheck errcheck lint vendor-status test-compile website website-test terraform-clean terraform-ipa remove-qa-container link-git-hooks

clean:
	rm -rf out
	go clean ./...

build: fmtcheck errcheck lint test
	mkdir -p out
	go build -o out/terraform-provider-$(PKG_NAME)

install: uninstall build
	@mkdir -p ~/.terraform.d/plugins
	@cp out/terraform-provider-$(PKG_NAME) ~/.terraform.d/plugins
	@echo "==> Installed provider at ~/.terraform.d/plugins/terraform-provider-$(PKG_NAME) ..."

uninstall:
	@rm -rf ~/.terraform.d/plugins/terraform-provider-$(PKG_NAME)

# Terraform Provider MongoDB: E2E test
TFDIR=examples/standalone-in-docker
terraform-ipa: terraform-clean clean install
	@echo "Initializing terraform, then applying the plan..."
	cd $(TFDIR); \
	terraform init; \
	terraform plan; \
	terraform apply -auto-approve

terraform-clean: remove-qa-container
	@echo "Destroying any existing resources and deleting TF state"
	-cd $(TFDIR); \
	terraform init; \
	terraform destroy -auto-approve; \
	rm -rf .terraform terraform.tfstate terraform.tfstate.backup *.log

remove-qa-container:
	@echo "Removing any existing QA containers..."
	$(MAKE) -C docker remove-all-containers

# GIT hooks
link-git-hooks:
	@echo "==> Installing all git hooks..."
	find .git/hooks -type l -exec rm {} \;
	find .githooks -type f -exec ln -sf ../../{} .git/hooks/ \;
