// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	space_parameter_name  = "space_param"
	space_sensitive_param = "space_sensitive_param"
	expected_value        = "value"
	space_name            = "TorqueTerraformProvider"
)

func TestSpaceParamResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Can't create account level tag with possible values
				Config: providerConfig + fmt.Sprintf(`
				data "torque_space_parameter" "space_param" {
					space_name            ="%s"
					name = "%s"
				}
				`, space_name, space_parameter_name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.torque_space_parameter.space_param", "name", space_parameter_name),
					resource.TestCheckResourceAttr("data.torque_space_parameter.space_param", "value", expected_value),
					resource.TestCheckResourceAttr("data.torque_space_parameter.space_param", "sensitive", "false"),
					resource.TestCheckResourceAttr("data.torque_space_parameter.space_param", "description", "description"),
				),
			},
			{
				// Can't create account level tag with possible values
				Config: providerConfig + fmt.Sprintf(`
				data "torque_space_parameter" "space_param" {
					space_name            ="%s"
					name = "%s"
				}
				`, space_name, space_sensitive_param),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.torque_space_parameter.space_param", "name", space_sensitive_param),
					resource.TestCheckNoResourceAttr("data.torque_space_parameter.space_param", "value"),
					resource.TestCheckResourceAttr("data.torque_space_parameter.space_param", "sensitive", "true"),
					resource.TestCheckResourceAttr("data.torque_space_parameter.space_param", "description", "description"),
				),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
				data "torque_space_parameter" "other_param" {
					space_name            ="%s"
					name = "some_non_exiting_param"
				}
				`, space_name),
				ExpectError: regexp.MustCompile(`Unable to Read Torque Space Parameter. Parameter was not found in space`),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
