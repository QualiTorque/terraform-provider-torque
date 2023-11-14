terraform {
  required_providers {
    torque = {
      source = "QualiTorque/torque"
    }
  }
}

provider "torque" {}

resource "torque_introspection_resource" "example" {
    display_name = "My Resource"
    image = "https://portal.qtorque.io/static/media/networking.dc1b7fb73182de0136d059a1eb00af4f.svg"
    introspection_data = {size = "large", mode = "party"}
}
