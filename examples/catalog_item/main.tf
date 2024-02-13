terraform {
  required_providers {
    torque = {
      source = "qualitorque/torque"
    }
  }
}

provider "torque" {
  host  = "https://portal.qtorque.io/"
  space = var.torque_space
  token = var.torque_token
}

resource "torque_catalog_item" "catalog_item" {
  space_name      = var.space_name
  blueprint_name  = var.blueprint_name
  repository_name = var.repository_name
}