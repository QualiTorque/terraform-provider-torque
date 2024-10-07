terraform {
  required_providers {
    torque = {
      source = "qualitorque/torque"
    }
  }
}

provider "torque" {
  # host  = "https://portal.qtorque.io/"
  # space = "api_space"
  # token = "111111111111"
}

resource "torque_space_git_credentials" "creds" {
  space_name  = "my_space"
  name        = "my_tf_creds"
  description = "desc"
  token       = "token"
  type        = "github"
}
