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

resource "torque_s3_object_content_input_source" "app-versions" {
  name                          = "app-versions"
  description                   = "description"
  specific_spaces               = ["space"]
  all_spaces                    = false
  bucket_name                   = "s3-bucket-name"
  credential_name               = "credentials"
  filter_pattern_overridable    = true
  filter_pattern                = "regex"
  json_path                     = "$.kubernetes.helm_release[*].name"
  object_key                    = "artifacts/builds/app/versions.json"
  object_key_overridable        = true
  json_path_overridable         = false
  display_json_path             = "$.kubernetes.helm_release[*].display_name"
  display_json_path_overridable = false
}
