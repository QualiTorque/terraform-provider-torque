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

resource "torque_account" "name" {
  parent_account = "torque_parent_account"
  account_name   = "my_sub_account"
  password       = "password"
  company        = "my company"
}
