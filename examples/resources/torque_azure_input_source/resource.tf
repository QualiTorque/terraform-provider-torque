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

resource "torque_azure_blob_object_input_source" "az_blob" {
  name        = "az-input-source"
  description = "az-input-source"
  # specific_spaces            = ["Space"]
  all_spaces                 = true
  storage_account_name       = "storage-account"
  container_name             = "container"
  credential_name            = "azure-creds"
  filter_pattern_overridable = true
  filter_pattern             = "regex"
  path_prefix                = "prefix"
  path_prefix_overridable    = false
}

resource "torque_azure_blob_object_content_input_source" "az_blob" {
  name        = "az-blob-object-content"
  description = "az-blob-object-content"
  # specific_spaces            = ["Space"]
  all_spaces                    = true
  storage_account_name          = "storage-account"
  container_name                = "container"
  blob_name                     = "blob"
  credential_name               = "azure-creds"
  filter_pattern_overridable    = true
  filter_pattern                = "regex"
  json_path                     = "$.kubernetes.helm_release[*].name"
  json_path_overridable         = false
  display_json_path             = "$.kubernetes.helm_release[*].display_name"
  display_json_path_overridable = false
}
