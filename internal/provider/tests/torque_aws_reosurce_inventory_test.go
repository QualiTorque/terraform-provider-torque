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
	name                           = "aws"
	description                    = "description"
	account_number                 = "11111111111"
	access_key                     = "key"
	secret_key                     = "secret"
	view_arn                       = "arn:aws:iam::123456789012:user/JohnDoe"
	new_name                       = "new_aws"
	new_description                = "new_description"
	new_account_number             = "22222222222"
	new_access_key                 = "new_key"
	new_secret_key                 = "new_secret"
	new_view_arn                   = "arn:aws:iam::123456789012:user/JohnnyDoe"
	invalid_view_arn               = "invalid_arn"
	invalid_view_arn_error_message = "must be a valid ARN format"
)

var unique_name = name + "_" + index
var new_unique_name = new_name + "_" + index

func TestTorqueAwsResourceInventory(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_aws_resource_inventory" "aws" {
	                name           = "%s"
					description    = "%s"
					account_number = "%s"
					access_key     = "%s"
					secret_key     = "%s"
					view_arn       = "%s"
				}
				`, unique_name, description, account_number, access_key, secret_key, view_arn),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_aws_resource_inventory.aws", "name", unique_name),
					resource.TestCheckResourceAttr("torque_aws_resource_inventory.aws", "description", description),
					resource.TestCheckResourceAttr("torque_aws_resource_inventory.aws", "account_number", account_number),
					resource.TestCheckResourceAttr("torque_aws_resource_inventory.aws", "access_key", access_key),
					resource.TestCheckResourceAttr("torque_aws_resource_inventory.aws", "secret_key", secret_key),
					resource.TestCheckResourceAttr("torque_aws_resource_inventory.aws", "view_arn", view_arn),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_aws_resource_inventory" "aws" {
	                name           = "%s"
					description    = "%s"
					account_number = "%s"
					access_key     = "%s"
					secret_key     = "%s"
					view_arn       = "%s"
				}
				`, new_unique_name, new_description, new_account_number, new_access_key, new_secret_key, new_view_arn),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_aws_resource_inventory.aws", "name", new_unique_name),
					resource.TestCheckResourceAttr("torque_aws_resource_inventory.aws", "description", new_description),
					resource.TestCheckResourceAttr("torque_aws_resource_inventory.aws", "account_number", new_account_number),
					resource.TestCheckResourceAttr("torque_aws_resource_inventory.aws", "access_key", new_access_key),
					resource.TestCheckResourceAttr("torque_aws_resource_inventory.aws", "secret_key", new_secret_key),
					resource.TestCheckResourceAttr("torque_aws_resource_inventory.aws", "view_arn", new_view_arn),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})

}

func TestRecreatedResourceInventory(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_aws_resource_inventory" "aws" {
	                name           = "%s"
					description    = "%s"
					account_number = "%s"
					access_key     = "%s"
					secret_key     = "%s"
					view_arn       = "%s"
				}
				`, unique_name, description, account_number, access_key, secret_key, invalid_view_arn),
				ExpectError: regexp.MustCompile(invalid_view_arn_error_message),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
