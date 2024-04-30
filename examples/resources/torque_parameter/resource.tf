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

resource "torque_parameter" "new_parameter" {
  name        = "parameter_name"
  value       = "parameter_value"
  sensitive   = true
  description = "parameter_description"
}
