terraform {
  required_providers {
    torque = {
      source = "qualitorque/torque"
    }
  }
}

provider "torque" {
  host  = ""
  space = ""
  token = ""
}

resource "torque_space_tag_value_association" "dev_space" {
  space_name = "dvir"
  tag_name   = "test"
  tag_value  = "Dev"
}
