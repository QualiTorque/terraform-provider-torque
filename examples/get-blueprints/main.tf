terraform {
  required_providers {
    torque = {
      source = "qualitorque/torque"
    }
  }
}

provider "torque" {
  host  = "https://portal.qtorque.io/"
  space = "01-Development"
  token = ""
}

data "torque_space_blueprints" "blueprints" {
  space_name = "01-Development"
  //filter by repo
}

