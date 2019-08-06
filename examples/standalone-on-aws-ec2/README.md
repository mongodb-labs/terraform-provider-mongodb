# Example Terraform HCL file for AWS EC2

- 1.Check out the latest available [MongoDB Terraform Provider](https://github.com/mongodb-labs/terraform-provider-mongodb/releases) 
  and then [build the provider plugin](https://github.com/mongodb-labs/terraform-provider-mongodb#setting-up-the-development-environment).

- 2.Install it
  ```bash
  make install # will copy it to `~/.terraform.d/plugins`
  ```

- 3.Prepare the required AWS dependencies:

    - A. a *security group* which should allow the following inbound traffic:
        - external access to the Ops Manager port
        - external access to the SSH port
        - access to the MongoDB port from within the same security group 
    
        Example:
        
        | Type            | Protocol | Port Range | Source    | Description         |
        | --------------- | -------- | ---------- | --------- | ------------------- |
        | Custom TCP Rule | TCP      | 9080       | 0.0.0.0/0 | Ops Manager         |
        | SSH             | TCP      | 22         | 0.0.0.0/0 | SSH                 |
        | Custom TCP Rule | TCP      | 27017      | sg-...    | This security group |

  - B. a *subnet* which automatically assigns an IPv4 address and allows inbound/outbound traffic.

- 4.Export the required variables
    - security_group_id
    - subnet_id

```bash
export TF_VAR_security_group_id=...
export TF_VAR_subnet_id=...
```

- 5.Run terraform

```
cd examples/standalone-on-aws-ec2 # this directory
make apply # which in turn will run terraform init, plan, and apply

# if you want to clean a previously existing plan, run:
make clean
```

### Connecting to the host

You can simply run `make ssh` to connect to the host.

If you need to export the generated private key, you can do so with `make private-key`, which will save it to a temporary file.
