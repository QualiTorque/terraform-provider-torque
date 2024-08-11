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
  scope           = "tag_scope"
  description     = "tag_description"
  possible_values = ["valu1", "value2"]
}