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

data "torque_space_blueprint" "blueprint" {
  name       = "blueprint"
  space_name = "target_space"
}
