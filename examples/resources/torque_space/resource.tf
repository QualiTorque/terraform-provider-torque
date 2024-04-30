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

resource "torque_space" "new_space" {
  space_name = "space"
  color      = "space_color"
  icon       = "space_icon"
}
