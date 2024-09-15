// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAWSCostTargetResource(t *testing.T) {
	var version = os.Getenv("VERSION")
	var minorVresion = strings.Split((version), ".")
	var index = minorVresion[1]
	new_aws_cost_target_name := fmt.Sprintf("new_aws_cost_target_%s", index)
	aws_cost_target_name := fmt.Sprintf("aws_cost_target_%s", index)
	random_account_num := fmt.Sprint(acctest.RandIntRange(100000000000, 999999999999))
	aws_cost_target_role_arn := fmt.Sprintf("arn:aws:iam::%s:role/role", random_account_num)
	new_aws_cost_target_role_arn := fmt.Sprintf("arn:aws:iam::%s:role/newrole", random_account_num)
	const (
		aws_cost_target_external_id     = "aws_cost_target_external_id"
		new_aws_cost_target_external_id = "aws_cost_target_external_id2"
		invalid_role_arn                = "some_arn"
	)
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
				ExpectError: regexp.MustCompile("Error Updating AWS Cost Target"),
			},
		},
	})
}
