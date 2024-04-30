terraform {
  required_providers {
    torque = {
      source = "qualitorque/torque"
    }
  }
}

provider "torque" {
  host  = "https://portal.qtorque.io/"
  space = var.torque_space
  token = var.torque_token
}


resource "torque_aws_cost_target" "cost_target" {
  name        = "cost_target_name"
  role_arn    = "arn:partition:service:region:account-id:resource-id"
  external_id = "external_id"
}