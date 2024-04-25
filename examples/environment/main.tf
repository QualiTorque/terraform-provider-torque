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

resource "torque_environment" "name" {
  environment_name = "blabla"
  blueprint_name   = "Hello World"
  duration         = "PT2H"
  space            = "Workshop"
  inputs = {
    "name" : "amir",
  }
  # owner_email = "amir.r@quali.com"
  automation  = true
  description = "Created by Terraform"
}
