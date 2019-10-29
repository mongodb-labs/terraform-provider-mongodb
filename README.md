MongoDB Terraform Provider
===================================

A [Terraform Provider](https://www.terraform.io/docs/providers/index.html) for MongoDB Cloud resources.

**This project is currently in development and is not yet ready for production use.**

This repository is based on [terraform-provider-scaffolding](https://github.com/terraform-providers/terraform-provider-scaffolding) 
and is licensed under the terms of the [Mozilla Public License Version 2.0](https://www.mozilla.org/en-US/MPL/2.0/).


#### Why write a terraform provider, when the same could be accomplished with bash code executed with the remote-exec provisioner?

As per [this page](https://www.terraform.io/docs/configuration/resources.html#provisioner-and-connection-resource-provisioners): 
> Provisioning steps should be used sparingly, since they represent non-declarative actions taken during the creation of a resource and so Terraform is not able to model changes to them as it can for the declarative portions of the Terraform language.

Managing MongoDB processes entails a number of complex scenarios, which are inefficient when handled via shell code.

The following actions are examples which, using a traditional (provisioner-based) approach, would result in the whole instance (AMI, VM, container, etc.) being shut-down and recreated from scratch:
- reconfiguring a replica set
- upgrading MongoDB to a newer version
- configuring MongoDB nodes to use SSL or require clients to use TLS/SSL certificates
- upgrading Ops Manager


### Desired feature set for release v0.1.0

- [x] Terraform Resource: download and install an unmanaged MongoD standalone
- [x] Terraform Resource: download and Install Ops Manager
- [x] Terraform: use version 0.12.x syntax
- [x] Ops Manager: register the first user (global owner)
- [x] Ops Manager: Create an agent key via the API
- [x] Terraform Resource: Install and configure the Automation Agent
- [ ] Terraform Resource: Enable Monitoring
- [ ] Terraform Resource: Deploy a new MongoD standalone (managed by Ops Manager)
- [ ] Terraform HCL: revisit the resource schema
- [ ] Terraform HCL: validate the resource inputs
- [ ] Code: investigate "one resource per go package" codebase re-organization
- [ ] Code quality: write [terraform plugin acceptance tests](https://www.terraform.io/docs/extend/testing/index.html)
- [ ] Code quality: write unit-tests


### Feature Backlog

- [ ] Terraform Resource: configure a MongoDB replica-set
- [ ] Terraform Resource: deploy a highly-available Ops Manager set-up
- [ ] Terraform Resource: configure unmanaged MongoDB with SSL
- [ ] Terraform Resource: install and configure Ops Manager Backup Daemon(s)
- [ ] Terraform Resource: handle Ops Manager upgrades / rolling upgrades
- [ ] Terraform Resource: gen-key generator and support importing existing keys
- [ ] Terraform Resource: deploy a managed replica set via Ops Manager Automation
- [ ] Terraform Resource: deploy a managed sharded cluster via Ops Manager Automation
- [ ] Terraform Resource: enable SSL for databases managed by Ops Manager Automation
- [ ] Terraform Data source: find MongoD binaries using a _version manifest_
- [ ] Terraform Data source: find the desired Ops Manager version using the _release archive_ 
- [ ] Terraform Resource: create an Organization in Ops Manager
- [ ] Terraform Resource: create a Project in Ops Manager
- [ ] Terraform Resource: create a User in Ops Manager
- [ ] Terraform Resource: investigate the difficulty of importing unmanaged resources (MongoD/Ops Manager/Backup Daemon installs)
- [ ] Terraform Resource: import managed databases resources (standalones, replica sets, sharded clusters)
- [ ] Terraform Resource: move a managed resource from Ops Manager to Cloud Manager (with monitoring downtime, no backups)
- [ ] Terraform Resource: move a managed resource from Ops Manager to Cloud Manager (fully online, including backups)
- [ ] Code quality: integrate with a code coverage service
- [ ] Terraform: write a test which checks that using a bastion host works as expected
- [ ] Terraform: manage MMS, AppDB, and AA as system services (when hosts restart, so do these)
- [ ] Terraform Resource: manage MMS, AppDB, and AA as system services (when hosts restart, so do these, e.g.: SystemD; user permissions)
- [ ] Terraform Data source: Automation, Monitoring, and Backup agent versions (to allow the provider to upgrade agents)

NOTE: the list above has not been prioritized. Please open a [New Issue](https://github.com/mongodb-labs/terraform-provider-mongodb/issues/new)
if you'd like to discuss them or submit new ideas.


### Setting up the development environment

Pull requests are always welcome! Please read our [contributor guide](./CONTRIB.md) before starting any work.  

The steps below should help you get started.  They have been tested on MacOS, but should work on Linux systems as well (with minor adaptations.)

1. Install GO
```
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
```

Ensure `$GOROOT/bin` is in your path.

2. Install the following tools

```
# GoLint
go get -u golang.org/x/lint/golint

# Golangci-lint
curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(go env GOPATH)/bin v1.21.0
```

3. Install the git hooks, to automatically fix linting issues and flag any errors 

`make link-git-hooks`

4. Download and install [Docker](https://hub.docker.com/editions/community/docker-ce-desktop-mac)

5. Download and install [GoLand](https://www.jetbrains.com/go/nextversion/)

NOTE: Goland is not a requirement; you could also use [Visual Studio Code](https://code.visualstudio.com/) or any other editor.

6. Install the following plugins
- Hashicorp Terraform HCL
- Makefile support
- Save actions
- Docker
- then restart GoLand
- and configure _Save Actions_ to:
  - reformat your code 
  - and to optimize imports

7. Import the project

- `cd /path/to/git && git clone ...`
- File -> New -> Project
- configure the path to an existing location (Goland will ask you later if it should import the existing project)
- Check that the project is using go modules (Preferences -> Go -> Go Modules -> Enable go modules integrations).
- If GoLang cannot resolve dependencies, enable go dep support (Preferences -> Go -> Dep -> Enable).

If the import fails, or you see errors related to `terraform-provider-mongodb.iml`, try the following:
- `rm -f .idea/terraform-provider-mongodb.iml`, try reattempt the import
- or if the above fails, `rm -rf .idea`, then reattempt the import.

If you're deleted `.idea`, you will need to reconfigure _File Watchers_ and the _Save Actions_ plugin (for the latter, see above.)

Open _File Watchers_ and add the following to the current project:
  - go fmt
  - goimports
  - terraform fmt
  - custom: golint (see [this post](https://github.com/vmware/dispatch/wiki/Configure-GoLand-with-golint))

8. Run an E2E test

**Prerequisite:**

The following command will create a Docker container to which you can connect to via SSH, using your SSH RSA key: `~/.ssh/id_rsa.pub`.

```bash
cd docker
make build-image
cd ..
```

Before going any further, ensure it is loaded by your SSH agent (`ssh-add`), or the next command will fail to connect. 

```bash
# E2E test
make terraform-ipa
# will run terraform init, plan, and apply, using a local Docker container
```

NOTE: to speed up local development, we cache the required MongoDB artifacts (MongoD and Ops Manager).  
While you can replace the `binary` links in `./examples/standalone-in-docker/` to their internet-based locations, 
it is much more efficient to download them once and serve them locally with [Miniserve](https://formulae.brew.sh/formula/miniserve)
or a similar static file web-server.  Once you've downloaded the two binaries, go to the download directory and run `miniserve -p 9000 .`


### Debugging terraform tests

```
make debugacc TESTARGS=-test.run=TestMongoDB_unit
```


### Installing the provider

**Make terraform aware of this plugin:**

```
make install
```

**Alternatively, use `.terraformrc`:**

```
cat >> ~/.terraformrc <<EOF
providers {
    mongodb = "/path/to/terraform-provider-mongodb"
}
EOF
```
