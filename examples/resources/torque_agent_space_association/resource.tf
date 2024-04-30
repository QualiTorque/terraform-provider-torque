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

resource "torque_agent_space_association" "agent_association" {
  space_name      = "space"
  agent_name      = "agent"
  service_account = "service-account"
  namespace       = "default"
}
