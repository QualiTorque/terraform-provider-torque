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

data "torque_space_parameter" "torque_space_parameter" {
  space_name = "target_space"
  name       = "space_parameter"
}
