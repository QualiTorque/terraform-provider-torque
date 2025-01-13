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

resource "torque_space_label" "label" {
  space_name   = "target_space"
  name         = "k8s"
  color        = "blue"
  quick_filter = true
}
