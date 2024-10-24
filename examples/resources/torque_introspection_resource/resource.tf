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
  display_name = "Custom card"
  image        = "https://raw.githubusercontent.com/QualiTorque/Torque-Samples/refs/heads/main/instructions/res_images/Tetris_logo_small.png"
  introspection_data = {
    "data1" : "value1"
    "data2" : "value2"
    "data3" : "value3"
  }
  links = [
    {
      "icon" : "upload",
      "href" : "https://example1.com"
      "label" : "example1"
      "color" : "#0000ff"
    },
    {
      "icon" : "download",
      "href" : "https://example2.com"
      "label" : "example2"
    },
    {
      "icon" : "copy",
      "href" : "https://example3.com"
      "label" : "example3"
      "color" : "#ff0000"
    },
    {
      "icon" : "play",
      "href" : "https://example4.com"
      "label" : "example4"
    },
    {
      "icon" : "restart",
      "href" : "https://example5.com"
      "label" : "example5"
      "color" : "#00ff00"
  }]
}
