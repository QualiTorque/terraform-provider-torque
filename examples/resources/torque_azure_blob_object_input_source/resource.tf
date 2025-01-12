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

resource "torque_azure_blob_object_input_source" "az_blob" {
  name                       = "az-blob-object-input-source"
  description                = "description"
  all_spaces                 = true
  storage_account_name       = "storage-account"
  container_name             = "container"
  credential_name            = "azure-creds"
  filter_pattern_overridable = true
  filter_pattern             = "regex"
  path_prefix                = "prefix"
  path_prefix_overridable    = false
}

