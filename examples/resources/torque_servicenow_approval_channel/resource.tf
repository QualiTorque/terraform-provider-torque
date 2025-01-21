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

resource "torque_servicenow_approval_channel" "channel" {
  name        = "name"
  description = "description"
  base_url    = "base_url"
  username    = "username"
  password    = "password"
  headers     = "json"
  approver    = "approver@company.com"
}
