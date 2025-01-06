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

resource "torque_s3_object_input_source" "s35" {
  name                       = "s3-object-input-source"
  description                = "s3-bucket-input-source-description"
  specific_spaces            = ["space1"]
  all_spaces                 = true
  bucket_name                = "s3-bucket-name"
  credential_name            = "AWS Credentials"
  filter_pattern_overridable = false
  filter_pattern             = "regex"
  path_prefix                = "prefix"
  path_prefix_overridable    = true
}
