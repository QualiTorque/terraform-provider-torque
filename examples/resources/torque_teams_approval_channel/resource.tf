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

resource "torque_teams_approval_channel" "channel" {
  name            = "channel"
  description     = "description"
  webhook_address = "webhook"
  approvers       = ["approver@company.com"]
}
