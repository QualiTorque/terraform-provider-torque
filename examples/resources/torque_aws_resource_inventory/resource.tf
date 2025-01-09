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

resource "torque_aws_resource_inventory" "aws" {
  name           = "account"
  description    = "description"
  account_number = "123456789012"
  access_key     = "key"
  secret_key     = "secret"
  view_arn       = "arn:aws:iam::123456789012:user/JohnDoe"
}
