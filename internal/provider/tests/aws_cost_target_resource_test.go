// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	aws_cost_target_name        = "aws_cost_target"
	aws_cost_target_role_arn    = "arn:aws:iam::123456789012:role/marketingadminrole"
	aws_cost_target_external_id = "aws_cost_target_external_id"
	invalid_role_arn            = "some_arn"
)

func TestAWSCostTargetResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Can't create account level tag with possible values
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_aws_cost_target" "test" {
					name = "%s"
					role_arn = "%s"
					external_id = "%s"
				}
				`, aws_cost_target_name, aws_cost_target_role_arn, aws_cost_target_external_id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_aws_cost_target.test", "name", aws_cost_target_name),
					resource.TestCheckResourceAttr("torque_aws_cost_target.test", "role_arn", aws_cost_target_role_arn),
					resource.TestCheckResourceAttr("torque_aws_cost_target.test", "external_id", aws_cost_target_external_id),
				),
			},
			{
				// Can't create account level tag with possible values
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_aws_cost_target" "test" {
					name = "%s"
					role_arn = "%s"
					external_id = "%s"
				}
				`, aws_cost_target_name, invalid_role_arn, aws_cost_target_external_id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_aws_cost_target.test", "name", aws_cost_target_name),
					resource.TestCheckResourceAttr("torque_aws_cost_target.test", "role_arn", invalid_role_arn),
					resource.TestCheckResourceAttr("torque_aws_cost_target.test", "external_id", aws_cost_target_external_id),
				),
			},
			// {
			// 	// Can't create account level tag with possible values
			// 	Config: providerConfig + fmt.Sprintf(`
			// 	resource "torque_aws_cost_target" "test" {
			// 		name = "%s"
			// 	}
			// 	`, aws_cost_target_name),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("torque_aws_cost_target.test", "name", aws_cost_target_name),
			// 		resource.TestCheckResourceAttr("torque_aws_cost_target.test", "role_arn", aws_cost_target_role_arn),
			// 		resource.TestCheckResourceAttr("torque_aws_cost_target.test", "external_id", aws_cost_target_external_id),
			// 	),
			// },
			// Delete testing automatically occurs in TestCase
		},
	})
}
