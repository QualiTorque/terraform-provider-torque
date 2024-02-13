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

resource "torque_blueprint_tag_value_association" "tag_association" {
  space_name        = var.space_name
  repo_name       = var.blueprint_repo_name
  tag_name   = var.tag_name
  tag_value = var.tag_value
  blueprint_name = var.blueprint_name
}