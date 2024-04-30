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
resource "torque_introspection_resource" "example" {
  display_name       = "resource_name"
  image              = "resource_image"
  introspection_data = { "data1" : "value1" }
}
