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

resource "torque_environment_label_association" "env_label_association" {
  space_name     = "space"
  environment_id = "vFc0DgneI8wa"
  labels         = [{ "key" : "app", "value" : "artifactory" }]
}
