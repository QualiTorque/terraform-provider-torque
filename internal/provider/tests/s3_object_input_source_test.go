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

func TestS3ObjectInputSource(t *testing.T) {
	const (
		input_name      = "s3"
		description     = "s3 input test"
		credential_name = "TerraformCreds"
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
				resource "torque_s3_object_input_source" "s3_bucket" {
					name                       = "%s"
					description                = "%s"
					bucket_name                = "%s"
					credential_name            = "%s"
					filter_pattern_overridable = false
					filter_pattern             = "pattern"
					path_prefix                = "prefix"
					path_prefix_overridable    = true
				}
				`, unique_name, description, unique_bucket, credential_name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_s3_object_input_source.s3_bucket", "name", unique_name),
					resource.TestCheckResourceAttr("torque_s3_object_input_source.s3_bucket", "bucket_name", unique_bucket),
					resource.TestCheckResourceAttr("torque_s3_object_input_source.s3_bucket", "credential_name", "TerraformCreds"),
					resource.TestCheckResourceAttr("torque_s3_object_input_source.s3_bucket", "all_spaces", "true"),
					resource.TestCheckResourceAttr("torque_s3_object_input_source.s3_bucket", "filter_pattern", "pattern"),
					resource.TestCheckResourceAttr("torque_s3_object_input_source.s3_bucket", "filter_pattern_overridable", "false"),
					resource.TestCheckResourceAttr("torque_s3_object_input_source.s3_bucket", "path_prefix", "prefix"),
					resource.TestCheckResourceAttr("torque_s3_object_input_source.s3_bucket", "path_prefix_overridable", "true"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_s3_object_input_source" "s3_bucket" {
					name                       = "%s_new"
					description                = "my-s3-bucket-input-source"
					specific_spaces            = ["TorqueTerraformProvider"]
					bucket_name                = "%s_new"
					credential_name            = "TerraformCreds2"
					filter_pattern_overridable = true
					filter_pattern             = "new_pattern"
					path_prefix                = "new_prefix"
					path_prefix_overridable    = false
				}
				`, unique_name, unique_bucket),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_s3_object_input_source.s3_bucket", "name", unique_name+"_new"),
					resource.TestCheckResourceAttr("torque_s3_object_input_source.s3_bucket", "bucket_name", unique_bucket+"_new"),
					resource.TestCheckResourceAttr("torque_s3_object_input_source.s3_bucket", "credential_name", "TerraformCreds2"),
					resource.TestCheckResourceAttr("torque_s3_object_input_source.s3_bucket", "filter_pattern", "new_pattern"),
					resource.TestCheckResourceAttr("torque_s3_object_input_source.s3_bucket", "filter_pattern_overridable", "true"),
					resource.TestCheckResourceAttr("torque_s3_object_input_source.s3_bucket", "path_prefix", "new_prefix"),
					resource.TestCheckResourceAttr("torque_s3_object_input_source.s3_bucket", "path_prefix_overridable", "false"),
					resource.TestCheckResourceAttr("torque_s3_object_input_source.s3_bucket", "all_spaces", "false"),
					resource.TestCheckResourceAttr("torque_s3_object_input_source.s3_bucket", "specific_spaces.0", "TorqueTerraformProvider"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})

}

func TestInvalidInputSourceConfiguration(t *testing.T) {
	// spaceName := os.Getenv("TORQUE_SPACE")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				resource "torque_s3_object_input_source" "s3_bucket" {
					name                       = "name"
					description                = "description"
					bucket_name                = "bucket"
					credential_name            = "creds"
					filter_pattern_overridable = false
					filter_pattern             = "pattern"
					path_prefix                = "prefix"
					path_prefix_overridable    = true
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
