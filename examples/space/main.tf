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

resource "torque_space" "new_space" {
  space_name = var.space_name
  color      = var.space_color
  icon       = var.space_icon
}