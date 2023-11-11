terraform {
  required_providers {
    torque = {
      source = "quali.com/terraform/torque"
    }
  }
}

provider "torque" {}

resource "torque_introspection_resource" "example" {
    display_name = "My Resource"
    image = "https://cdn-icons-png.flaticon.com/512/882/882730.png"
    introspection_data = {size = "large", mode = "party"}
}