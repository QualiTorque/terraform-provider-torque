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

resource "torque_group" "new_group" {
  group_name     = "group_name"
  description    = "group_description"
  idp_identifier = "group_idp_identifier"
  users          = ["user1,user2"]
  account_role   = "group_account_role"
  space_roles = [{
    space_name = "space",
    space_role = "role"
  }]
}
