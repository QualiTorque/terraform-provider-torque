terraform {
  required_providers {
    torque = {
      source = "qualitorque/torque"
    }
  }
}

provider "torque" {
  host  = "https://portal.qtorque.io/"
  space = "api_space"
  token = "111111111111"
}

data "torque_environment_introspection" "env_introspection" {
  space_name = "target_space"
  id         = "JL4kgRgxT3Vo"
}
