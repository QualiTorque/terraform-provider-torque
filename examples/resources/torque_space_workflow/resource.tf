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

resource "torque_space_workflow" "workflow" {
  name            = "Space Day2 Operation"
  space_name      = "Space"
  repository_name = "repo"
}
