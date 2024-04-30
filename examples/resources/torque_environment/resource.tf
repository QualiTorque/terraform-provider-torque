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
  environment_name = var.environment_name
  blueprint_name   = var.blueprint_name
  duration         = var.duration
  space            = var.space
  owner_email      = var.owner_email
  inputs           = var.inputs
  collaborators    = var.collaborators
  automation       = var.automation
  description      = var.description
  tags             = var.tags
  workflows        = var.workflows
  blueprint_source = var.blueprint_source
}
