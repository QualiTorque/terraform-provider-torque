terraform {
  required_providers {
    torque = {
      source = "qualitorque/torque"
    }
  }
}

# provider "torque" {
#   host  = "https://portal.qtorque.io/"
#   space = "space"
#   token = "111111111111"
# }

resource "torque_aws_resource_inventory" "aws" {
  credentials = "creds"
  view_arn    = "view_arn"
}
