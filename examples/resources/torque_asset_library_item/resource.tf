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

resource "torque_asset_library_item" "library_item" {
  space_name      = "space"
  blueprint_name  = "blueprint"
  repository_name = "repository"
}
