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

data "torque_environment_introspection" "env_introspection" {
  space_name = var.torque_space
  id         = "JL4kgRgxT3Vo"
}
