# Quali's Torque Terraform Provider

The Terraform Provider for Quali's Torque is a plugin for Terraform that allows you to interact with Torque and control Torque behavior and presentation.

Learn more:

* Read more about [Quali's Torque](https://www.quali.com/torque/).

* Read more in the Torque [documentation](http://docs.qtorque.io/).

* Join the community [discussions](https://github.com/QualiTorque/qualitorque.github.io/).

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Quick Start

```terraform
terraform {
  required_providers {
    torque = {
      source = "QualiTorque/torque"
      version = "0.0.1"
    }
  }
}

provider "torque" {}

```


## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```
