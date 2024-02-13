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

resource "torque_parameter" "new_parameter" {
  name        = var.parameter_name
  value       = var.parameter_value
  sensitive   = var.parameter_sensitive
  description = var.parameter_description
}