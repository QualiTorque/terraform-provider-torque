// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestS3ObjectContentInputSource(t *testing.T) {
	const (
		input_name            = "s3"
		description           = "s3 input test"
		space                 = "TorqueTerraformProvider"
		credential_name       = "TerraformCreds"
		new_credential_name   = "TerraformCreds2"
		pattern               = "pattern"
		new_pattern           = "new_pattern"
		json_path             = "/"
		new_json_path         = "/new/path"
		display_json_path     = "/json/display/path"
		new_display_json_path = "/new/display/json/path"
		new_key               = "new_key"
		key                   = "key"
	)
	var version = os.Getenv("VERSION")
	var minorVresion = strings.Split((version), ".")
	var index = minorVresion[1]
	var unique_name = input_name + "_" + index
	var unique_bucket = input_name + "_bucket_" + index

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_s3_object_content_input_source" "s3_bucket" {
					name                          = "%s"
					description                   = "%s"
					bucket_name                   = "%s"
					credential_name               = "%s"
					filter_pattern_overridable    = false
					filter_pattern                = "%s"
					json_path                     = "%s"
					json_path_overridable         = false
					object_key                    = "%s"
					display_json_path             = "%s"
					display_json_path_overridable = true
				}
				`, unique_name, description, unique_bucket, credential_name, pattern, json_path, key, display_json_path),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_s3_object_content_input_source.s3_bucket", "name", unique_name),
					resource.TestCheckResourceAttr("torque_s3_object_content_input_source.s3_bucket", "bucket_name", unique_bucket),
					resource.TestCheckResourceAttr("torque_s3_object_content_input_source.s3_bucket", "credential_name", credential_name),
					resource.TestCheckResourceAttr("torque_s3_object_content_input_source.s3_bucket", "all_spaces", "true"),
					resource.TestCheckResourceAttr("torque_s3_object_content_input_source.s3_bucket", "filter_pattern", pattern),
					resource.TestCheckResourceAttr("torque_s3_object_content_input_source.s3_bucket", "filter_pattern_overridable", "false"),
					resource.TestCheckResourceAttr("torque_s3_object_content_input_source.s3_bucket", "object_key", key),
					resource.TestCheckResourceAttr("torque_s3_object_content_input_source.s3_bucket", "object_key_overridable", "false"),
					resource.TestCheckResourceAttr("torque_s3_object_content_input_source.s3_bucket", "json_path", json_path),
					resource.TestCheckResourceAttr("torque_s3_object_content_input_source.s3_bucket", "json_path_overridable", "false"),
					resource.TestCheckResourceAttr("torque_s3_object_content_input_source.s3_bucket", "display_json_path", display_json_path),
					resource.TestCheckResourceAttr("torque_s3_object_content_input_source.s3_bucket", "display_json_path_overridable", "true"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_s3_object_content_input_source" "s3_bucket" {
					name                          = "%s_new"
					description                   = "%s_new"
					bucket_name                   = "%s_new"
					credential_name               = "%s"
					filter_pattern_overridable    = true
					filter_pattern                = "%s"
					json_path                     = "%s"
					json_path_overridable         = true
					object_key                    = "%s"
					display_json_path             = "%s"
					display_json_path_overridable = false
					specific_spaces               = ["%s"]
				}
				`, unique_name, description, unique_bucket, new_credential_name, new_pattern, new_json_path, new_key, new_display_json_path, space),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_s3_object_content_input_source.s3_bucket", "name", unique_name+"_new"),
					resource.TestCheckResourceAttr("torque_s3_object_content_input_source.s3_bucket", "bucket_name", unique_bucket+"_new"),
					resource.TestCheckResourceAttr("torque_s3_object_content_input_source.s3_bucket", "description", description+"_new"),
					resource.TestCheckResourceAttr("torque_s3_object_content_input_source.s3_bucket", "credential_name", new_credential_name),
					resource.TestCheckResourceAttr("torque_s3_object_content_input_source.s3_bucket", "filter_pattern", new_pattern),
					resource.TestCheckResourceAttr("torque_s3_object_content_input_source.s3_bucket", "filter_pattern_overridable", "true"),
					resource.TestCheckResourceAttr("torque_s3_object_content_input_source.s3_bucket", "json_path", new_json_path),
					resource.TestCheckResourceAttr("torque_s3_object_content_input_source.s3_bucket", "json_path_overridable", "true"),
					resource.TestCheckResourceAttr("torque_s3_object_content_input_source.s3_bucket", "display_json_path", new_display_json_path),
					resource.TestCheckResourceAttr("torque_s3_object_content_input_source.s3_bucket", "display_json_path_overridable", "false"),
					resource.TestCheckResourceAttr("torque_s3_object_content_input_source.s3_bucket", "all_spaces", "false"),
					resource.TestCheckResourceAttr("torque_s3_object_content_input_source.s3_bucket", "specific_spaces.0", space),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})

}

func TestCantSpecifyAllSpacesWithSpecificSpacesInObjectContent(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				resource "torque_s3_object_content_input_source" "s3_bucket" {
					name                       = "name"
					description                = "description"
					bucket_name                = "bucket"
					credential_name            = "creds"
					filter_pattern_overridable = false
					filter_pattern             = "pattern"
					specific_spaces            = ["TorqueTerraformProvider"]
					all_spaces = true
				}
				`,
				ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
