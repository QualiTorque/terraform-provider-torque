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


resource "torque_workflow" "space_workflow" {
  name            = "Space Workflow"
  space_name      = "project-dev"
  repository_name = "repo"
}

resource "torque_workflow" "env_workflow" {
  name            = "Env Day2"
  space_name      = "project-dev"
  repository_name = "repo"
  self_service    = true
}
