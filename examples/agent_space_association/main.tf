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

resource "torque_agent_space_association" "agent_association" {
  space_name      = var.space_name
  agent_name      = var.agent_name
  service_account = var.service_account
  namespace       = var.namespace
}