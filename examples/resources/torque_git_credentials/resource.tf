terraform {
  required_providers {
    torque = {
      source = "qualitorque/torque"
    }
  }
}

provider "torque" {
  host  = "https://portal.qtorque.io/"
  space = "api_space"
  token = "111111111111"
}

resource "torque_git_credentials" "credentials" {
  name                = "credentials"
  description         = "description"
  token               = "token"
  type                = "github"
  allowed_space_names = ["Space1", "Space2"]
}
