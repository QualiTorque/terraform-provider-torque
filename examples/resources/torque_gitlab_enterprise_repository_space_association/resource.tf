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
  repository_name = "repository_name"
  repository_url  = "repository_url"
  token           = "token"
  branch          = "branch"
  credential_name = "credentials"
  use_all_agents  = false
  agents          = ["eks", "aks"]
}
