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

resource "torque_email_approval_channel" "channel" {
  name        = "channel"
  description = "description"
  approvers   = ["approver@company.com"]
}
