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

resource "torque_workflow" "workflow" {
  name            = "Day2 Operation"
  space_name      = "target_space"
  repository_name = "repo"
}
