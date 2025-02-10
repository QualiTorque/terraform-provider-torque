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
  space_name              = "space"
  blueprint_name          = "blueprint_name"
  repository_name         = "repository_name"
  self_service            = true
  display_name            = "display_name"
  max_duration            = "PT2H"
  default_duration        = "PT2H"
  default_extend          = "PT2H"
  max_active_environments = 10
  always_on               = false
  allow_scheduling        = true
  custom_icon             = "blueprint_icons/key"
  labels                  = ["label1", "label2"]
}
