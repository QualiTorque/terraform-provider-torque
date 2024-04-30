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

resource "torque_space_parameter" "new_parameter" {
  space_name  = "space"
  name        = "parameter_name"
  value       = "parameter_value"
  sensitive   = true
  description = "parameter_description"
}