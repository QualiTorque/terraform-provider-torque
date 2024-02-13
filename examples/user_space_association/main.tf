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

resource "torque_user_space_association" "user_association" {
  space_name = var.space_name
  user_email       = var.user_email
  user_role       = var.user_role
}