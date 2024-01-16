terraform {
  required_providers {
    torque = {
      source = "qualitorque/torque"
    }
  }
}

provider "torque" {
  host  = "https://portal.qtorque.io/"
  space = "03-Live"
  token = ""
}

resource "torque_blueprint_resource" "imported_blueprint" {
  name = "unknown"
}
