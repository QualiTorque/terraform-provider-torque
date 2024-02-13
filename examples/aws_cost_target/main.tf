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
  name        = var.cost_target_name
  role_arn    = var.role_arn
  external_id = var.external_id
}