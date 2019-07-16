# Configure providers
provider "docker" {
  host = "unix:///var/run/docker.sock"
}

# Create a container
resource "docker_container" "mdb0-0" {
  name  = "qa_ubuntu1804"
  image = "mongodb-terraform/qa_ubuntu1804"
  start = true
  ports {
    internal = 22
    external = 33022
  }
  ports {
    internal = 8080
    external = 33080
  }
}

# Deploy a MongoDB standalone
resource "mongodb_process" "mdb_standalone" {
  host {
    user     = "root"
    hostname = "127.0.0.1"
    port     = docker_container.mdb0-0.ports[0].external
  }

  mongod {
    # http://downloads.mongodb.org/linux/mongodb-linux-x86_64-ubuntu1804-4.0.10.tgz
    binary  = "http://localhost:9000/mongodb-linux-x86_64-ubuntu1804-4.0.10.tgz"
    bindip  = "0.0.0.0"
    port    = 27017
    workdir = "/opt/mongodb"
  }
}

# Generate an encryption key for Ops Manager
resource "random_string" "encryptionkey" {
  length  = 24
  special = true
}

# Deploy a single instance of Ops Manager
resource "mongodb_opsmanager" "opsman" {
  host {
    user     = "root"
    hostname = "127.0.0.1"
    port     = docker_container.mdb0-0.ports[0].external
  }

  opsmanager {
    # https://downloads.mongodb.com/on-prem-mms/deb/mongodb-mms_4.1.8.55362.20190620T1446Z-1_x86_64.deb
    binary              = "http://localhost:9000/mongodb-mms_4.1.8.55966.20190620T2143Z-1_x86_64.deb"
    workdir             = "/opt/mongodb"
    mongo_uri           = "mongodb://${mongodb_process.mdb_standalone.host.0.hostname}:${mongodb_process.mdb_standalone.mongod.0.port}/"
    encryption_key      = random_string.encryptionkey.result
    port                = 8080
    central_url         = "http://${mongodb_process.mdb_standalone.host.0.hostname}:8080"
    register_first_user = true

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
      "automation.agent.version"      = "10.2.0.5851-1"
    }
  }
}

resource "mongodb_automation_agent" "automation_agent" {
  host {
    user     = "root"
    hostname = "127.0.0.1"
    port     = docker_container.mdb0-0.ports[0].external
  }

  automation {
    binary     = "${mongodb_opsmanager.opsman.opsmanager[0].central_url}/download/agent/automation/mongodb-mms-automation-agent-10.2.0.5851-1.linux_x86_64.tar.gz"
    baseurl    = mongodb_opsmanager.opsman.opsmanager[0].central_url
    project_id = "fakeProjectId1"
    api_key    = "fakeApiKey"
    overrides = {
      serverPoolKey                 = "fakeServerPoolKeyShouldBeReplaced"
      sslMMSServerClientCertificate = "fakeCertificatePathSettingShouldBeAdded"
    }
  }
}
