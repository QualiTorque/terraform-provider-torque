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

resource "torque_tag" "new_tag" {
  name            = "tag_name"
  value           = "tag_value"
  scope           = "space"
  description     = "tag_description"
  possible_values = ["value1", "value2"]
}