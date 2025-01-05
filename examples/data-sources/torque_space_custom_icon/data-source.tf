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


data "torque_space_custom_icon" "icon" {
  space_name = "space"
  file_name  = "icon.svg"
}

output "key" {
  value = data.torque_space_custom_icon.icon.key
}
