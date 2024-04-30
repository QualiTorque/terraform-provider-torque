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


resource "torque_blueprint_tag_value_association" "tag_association" {
  space_name     = "space"
  repo_name      = "blueprint_repo_name"
  tag_name       = "tag_name"
  tag_value      = "tag_value"
  blueprint_name = "blueprint_name"
}