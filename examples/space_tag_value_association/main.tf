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

resource "torque_space_tag_value_association" "tag_association" {
  space_name = var.space_name
  tag_name       = var.tag_name
  tag_value       = var.tag_value
}