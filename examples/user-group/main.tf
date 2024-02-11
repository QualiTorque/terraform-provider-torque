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

resource "torque_group" "new_group" {
  group_name  = var.group_name
  description = var.group_description
  idp_identifier = var.group_idp_identifier
  users = var.group_users
  account_role =  var.group_account_role
  space_roles = var.group_space_roles
}