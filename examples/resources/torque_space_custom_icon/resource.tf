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

resource "torque_space_custom_icon" "icon" {
  space_name = "target_space"
  file_path  = "/path/to/svg/file.svg"
}
