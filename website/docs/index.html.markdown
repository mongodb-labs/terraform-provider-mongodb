---
layout: "mongodb"
page_title: "Provider: MongoDB"
sidebar_current: "docs-mongodb-index"
description: |-
    MongoDB is a cross-platform document-oriented database program. 
    The MongoDB Terraform provider allow for easy installation and management
    of MongoDB databases.
---

# MongoDB Provider

[MongoDB](https://mongodb.com/) is a cross-platform document-oriented database program. 

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure a MongoDB standalone
resource "mongodb_process" "standalone" {
  host {
    user     = "root"
    hostname = "127.0.0.1"
    port     = 22
  }

  mongod {
    binary  = "http://downloads.mongodb.org/linux/mongodb-linux-x86_64-ubuntu1804-4.0.10.tgz"
    bindip = "0.0.0.0"
    port   = 27017
  }
}
```