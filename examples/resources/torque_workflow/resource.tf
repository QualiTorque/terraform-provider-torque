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


resource "torque_workflow" "env_workflow" {
  name            = "Day2"
  space_name      = "Space"
  repository_name = "Repo"
  self_service    = true
  custom_icon     = "blueprint_icons/key"
}

