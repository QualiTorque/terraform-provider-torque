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

resource "torque_user_space_association" "user_association" {
  space_name = "space_name"
  user_email = "user_email"
  user_role  = "user_role"
}