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

resource "torque_tag" "new_tag" {
  name            = var.tag_name
  value           = var.tag_value
  scope           = var.tag_scope
  description     = var.tag_description
  possible_values = var.tag_possible_values
}