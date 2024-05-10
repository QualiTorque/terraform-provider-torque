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

resource "torque_codecommit_repository_space_association" "repository" {
  space           = "space"
  token           = "token"
  aws_region      = "eu-west-1"
  external_id     = "external_id"
  role_arn        = "arn:aws:iam::111111111111:role/role"
  git_username    = "CodeCommituser-at-111111111111"
  git_password    = "password"
  branch          = "main"
  repository_name = "codecommit-repo"
  repository_url  = "https://git-codecommit.eu-west-1.amazonaws.com/v1/repos/repo"
}
