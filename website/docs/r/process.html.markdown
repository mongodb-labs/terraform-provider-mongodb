---
layout: "mongodb"
page_title: "MongoDB: mongodb_process"
sidebar_current: "docs-mongodb-resource-process"
description: |-
    Create and manage a MongoDB process.
---

# mongodb\_process

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

## Argument Reference

TODO.

## Attributes Reference

TODO.

## Import

TODO.

```
$ terraform import ...
```
