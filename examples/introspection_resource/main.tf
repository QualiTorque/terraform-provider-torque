terraform {
  required_providers {
    torque = {
      source = "qualitorque/torque"
    }
  }
}

provider "torque" {}

resource "torque_introspection_resource" "example" {
  display_name       = var.resource_name
  image              = var.resource_image
  introspection_data = var.resource_data
}