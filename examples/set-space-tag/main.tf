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

resource "torque_tag_resource" "dev_space" {
  name  = "mytag"
  value = "mytagvalue"
  scope  = "space"
}
