// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAzureBlobObjectContentInputSourceResource(t *testing.T) {
	const (
		input_name               = "az_blob_content"
		description              = "az input test"
		space                    = "TorqueTerraformProvider"
		container_name           = "container"
		new_container_name       = "new_container"
		storage_account_name     = "sa"
		new_storage_account_name = "new_sa"
		credential_name          = "azure-creds"
		new_credential_name      = "__quali_azure__"
		blob_name                = "blob"
		new_blob_name            = "new_blob"
		pattern                  = "pattern"
		new_pattern              = "new_pattern"
		json_path                = "/"
		new_json_path            = "/new/path"
		path_prefix              = "prefix"
		new_path_prefix          = "new_prefix"
		display_json_path        = "/json/display/path"
		new_display_json_path    = "/new/display/json/path"
	)
	var unique_name = input_name + "_" + index
	var new_unique_name = "new" + "_" + input_name + "_" + index

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_azure_blob_object_content_input_source" "az_blob" {
					name        = "%s"
					description = "%s"
					specific_spaces            = ["%s"]
					// all_spaces                 = "false"
					storage_account_name       = "%s"
					container_name             = "%s"
					blob_name = "%s"
					credential_name            = "%s"
					filter_pattern_overridable = "true"
					filter_pattern             = "%s"
					json_path                     = "%s"
					json_path_overridable         = "false"
					display_json_path             = "%s"
					display_json_path_overridable = "false"
				}
				`, unique_name, description, space, storage_account_name, container_name, blob_name, credential_name, pattern, json_path, display_json_path),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "name", unique_name),
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "storage_account_name", storage_account_name),
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "credential_name", credential_name),
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "all_spaces", "false"),
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "filter_pattern", pattern),
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "filter_pattern_overridable", "true"),
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "container_name", container_name),
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "container_name_overridable", "false"),
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "blob_name", blob_name),
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "blob_name_overridable", "false"),
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "filter_pattern_overridable", "true"),
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "json_path", json_path),
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "json_path_overridable", "false"),
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "display_json_path", display_json_path),
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "display_json_path_overridable", "false"),
				),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_azure_blob_object_content_input_source" "az_blob" {
					name        = "%s"
					description = "%s"
					all_spaces                 = "true"
					storage_account_name       = "%s"
					container_name             = "%s"
					container_name_overridable = "true"
					blob_name = "%s"
					blob_name_overridable = "true"
					credential_name            = "%s"
					filter_pattern_overridable = "false"
					filter_pattern             = "%s"
					json_path                     = "%s"
					json_path_overridable         = "false"
					display_json_path             = "%s"
					display_json_path_overridable = "false"
				}
				`, new_unique_name, new_description, new_storage_account_name, new_container_name, new_blob_name, new_credential_name, new_pattern, new_json_path, new_display_json_path),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "name", new_unique_name),
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "storage_account_name", new_storage_account_name),
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "credential_name", new_credential_name),
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "all_spaces", "true"),
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "filter_pattern", new_pattern),
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "filter_pattern_overridable", "false"),
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "container_name", new_container_name),
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "container_name_overridable", "true"),
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "blob_name", new_blob_name),
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "blob_name_overridable", "true"),
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "json_path", new_json_path),
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "json_path_overridable", "false"),
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "display_json_path", new_display_json_path),
					resource.TestCheckResourceAttr("torque_azure_blob_object_content_input_source.az_blob", "display_json_path_overridable", "false"),
				),
			},
		},
	})
}
