terraform {
  required_providers {
    torque = {
      source = "qualitorque/torque"
    }
  }
}

provider "torque" {
  host  = ""
  space = ""
  token = ""
}

data "torque_space_blueprints" "blueprints" {
  space_name = ""
  //filter by repo
}

