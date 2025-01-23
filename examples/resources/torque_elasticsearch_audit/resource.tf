terraform {
  required_providers {
    torque = {
      source = "qualitorque/torque"
    }
  }
}

# provider "torque" {
#   host  = "https://portal.qtorque.io/"
#   space = "space"
#   token = "111111111111"
# }

resource "torque_elasticsearch_audit" "audit" {
  url      = "https://elastic:9000"
  username = "user"
  password = "password"
}
