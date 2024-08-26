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
	aws_cost_target_name            = "aws_cost_target"
	aws_cost_target_role_arn        = "arn:aws:iam::123456789012:role/marketingadminrole"
	new_aws_cost_target_role_arn    = "arn:aws:iam::12345678111:role/marketingadminrole"
	aws_cost_target_external_id     = "aws_cost_target_external_id"
	new_aws_cost_target_external_id = "aws_cost_target_external_id2"
	invalid_role_arn                = "some_arn"
	new_aws_cost_target_name        = "new_aws_cost_target"
)

func TestAWSCostTargetResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
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
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_aws_cost_target" "test" {
					name = "%s"
					role_arn = "%s"
					external_id = "%s"
				}
				`, new_aws_cost_target_name, new_aws_cost_target_role_arn, new_aws_cost_target_external_id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_aws_cost_target.test", "name", new_aws_cost_target_name),
					resource.TestCheckResourceAttr("torque_aws_cost_target.test", "role_arn", new_aws_cost_target_role_arn),
					resource.TestCheckResourceAttr("torque_aws_cost_target.test", "external_id", new_aws_cost_target_external_id),
				),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_aws_cost_target" "test" {
					name = "%s"
					role_arn = "%s"
					external_id = "%s"
				}
				`, new_aws_cost_target_name, invalid_role_arn, new_aws_cost_target_external_id),
				ExpectError: regexp.MustCompile("Unable to create AWS cost collection target"),
			},
		},
	})
}
