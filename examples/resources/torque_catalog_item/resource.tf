terraform {
  required_providers {
    torque = {
      source = "qualitorque/torque"
    }
  }
}

provider "torque" {
  host  = "https://portal.qtorque.io/"
  space = "space"
  token = "111111111111"
}

resource "torque_catalog_item" "catalog_item" {
  space_name      = "space"
  blueprint_name  = "blueprint_name"
  repository_name = "repository_name"
}
