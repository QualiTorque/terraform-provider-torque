terraform {
  required_providers {
    torque = {
      source = "qualitorque/torque"
    }
  }
}

provider "torque" {
  host  = var.torque_host
  space = var.torque_space
  token = var.torque_token
}

resource "torque_repository_space_association" "repository" {
  space_name      = var.space_name
  repository_url  = var.repository_url
  access_token    = var.access_token
  repository_type = var.repository_type
  branch          = var.branch
  repository_name = var.repository_name
}
