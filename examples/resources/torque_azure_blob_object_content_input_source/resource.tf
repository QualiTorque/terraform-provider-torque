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

resource "torque_azure_blob_object_content_input_source" "az_blob" {
  name                          = "az-blob-object-content"
  description                   = "az-blob-object-content"
  specific_spaces               = ["space"]
  storage_account_name          = "storage-account"
  container_name                = "container"
  blob_name                     = "blob"
  blob_name_overridable         = true
  credential_name               = "azure-creds"
  filter_pattern_overridable    = true
  filter_pattern                = "regex"
  json_path                     = "$.kubernetes.helm_release[*].name"
  json_path_overridable         = true
  display_json_path             = "$.kubernetes.helm_release[*].display_name"
  display_json_path_overridable = false
}
