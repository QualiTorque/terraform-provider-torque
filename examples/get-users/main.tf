terraform {
  required_providers {
    torque = {
      source = "qualitorque/torque"
    }
  }
}

provider "torque" {
  host  = "https://review2.qualilabs.net/"
  space = "Shira"
  token = ""
}

data "torque_user" "edu" {
  user_email = "tomer.a@quali.com"
}

output "edu_coffees" {
  value = data.torque_user.edu
}