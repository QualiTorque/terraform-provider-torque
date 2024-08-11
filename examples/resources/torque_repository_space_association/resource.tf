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

resource "torque_repository_space_association" "repository" {
  space_name      = "space_name"
  repository_url  = "repository_url"
  access_token    = "access_token"
  repository_type = "repository_type"
  branch          = "branch"
  repository_name = "repository_name"
}
