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

data "torque_space_blueprints" "torque_space_blueprints" {
  space_name                = "name"
  filter_by_repository_name = "repo"
}

