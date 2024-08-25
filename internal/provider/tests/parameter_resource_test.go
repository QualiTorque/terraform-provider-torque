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
	account_parameter_name  = "tf_provider_test_param"
	account_sensitive_param = "tf_provider_test_sensitive_param"
	// expected_value          = "value"
	// space_name              = "TorqueTerraformProvider"
)

func TestAccountParameterDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Can't create account level tag with possible values
				Config: providerConfig + fmt.Sprintf(`
				data "torque_parameter" "param" {
					name = "%s"
				}
				`, account_parameter_name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.torque_parameter.param", "name", account_parameter_name),
					resource.TestCheckResourceAttr("data.torque_parameter.param", "value", expected_value),
					resource.TestCheckResourceAttr("data.torque_parameter.param", "sensitive", "false"),
					resource.TestCheckResourceAttr("data.torque_parameter.param", "description", "description"),
				),
			},
			{
				// Can't create account level tag with possible values
				Config: providerConfig + fmt.Sprintf(`
				data "torque_parameter" "param" {
					name = "%s"
				}
				`, account_sensitive_param),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.torque_parameter.param", "name", account_sensitive_param),
					resource.TestCheckNoResourceAttr("data.torque_parameter.param", "value"),
					resource.TestCheckResourceAttr("data.torque_parameter.param", "sensitive", "true"),
					resource.TestCheckResourceAttr("data.torque_parameter.param", "description", "description"),
				),
			},
			{
				Config: providerConfig + `
				data "torque_parameter" "other_param" {
					name = "some_non_exiting_param"
				}
				`,
				ExpectError: regexp.MustCompile(`Unable to Read Torque Account Parameter`),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
