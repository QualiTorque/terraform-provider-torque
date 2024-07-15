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

resource "torque_gitlab_enterprise_repository_space_association" "repository" {
  space_name      = "space_name"
  base_url        = "base_url"
  repository_name = "repository_name"
  repository_url  = "repository_url"
  token           = "token"
  branch          = "branch"
}
