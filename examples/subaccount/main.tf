terraform {
  required_providers {
    torque = {
      source = "qualitorque/torque"
    }
  }
}

provider "torque" {
  host  = "https://review1.qualilabs.net/"
  space = var.torque_space
  token = var.torque_token
}

resource "torque_account" "name" {
  parent_account = "trial-3ba8f8b0"
  account_name   = "amir-tf1"
  password       = "Zubur123!"
  company        = "Quali1"
}
