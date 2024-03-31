terraform {
  required_providers {
    torque = {
      source = "qualitorque/torque"
    }
  }
}

provider "torque" {
  host  = "https://portal.qtorque.io"
  space = var.torque_space
  token = var.torque_token
}

resource "torque_account" "name" {
  parent_account = var.parent_account
  account_name   = var.account_name
  password       = var.password
  company        = var.company
}
