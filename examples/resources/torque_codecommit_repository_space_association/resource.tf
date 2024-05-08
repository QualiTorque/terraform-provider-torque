terraform {
  required_providers {
    torque = {
      source = "qualitorque/torque"
    }
  }
}

provider "torque" {
  host  = var.host
  space = var.space
  token = var.token
}

resource "torque_codecommit_repository_space_association" "repository" {
  space_name      = var.space
  aws_region      = var.aws_region
  role_arn        = var.role_arn
  external_id     = var.external_id
  username        = var.username
  password        = var.password
  branch          = var.branch
  repository_name = var.repository_name
  repository_url  = var.repository_url
}
