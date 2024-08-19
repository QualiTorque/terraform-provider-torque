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

resource "torque_space_label_association" "label_association" {
  space_name      = "target_space"
  blueprint_name  = "blueprint"
  repository_name = "repo"
  labels          = ["k8s", "aws"]
}
