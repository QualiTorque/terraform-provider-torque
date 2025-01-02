terraform {
  required_providers {
    torque = {
      source = "qualitorque/torque"
    }
  }
}

provider "torque" {
  host  = "https://portal.qtorque.io/"
  space = "api_space"
  token = "111111111111"
}

resource "torque_input_source" "input_source" {
  name            = "s3_input_source"
  description     = "description"
  type            = "s3"
  all_spaces      = false
  specific_spaces = ["space1", "space2"]
}
