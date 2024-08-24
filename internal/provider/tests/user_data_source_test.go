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
	user_email               = "terraformtester@quali.com"
	first_name               = "Developer"
	last_name                = "McDeveloper"
	display_first_name       = "Developer"
	display_last_name        = "McDeveloper"
	user_type                = "REGULAR"
	join_date                = "2024-08-23T11:20:56.827262"
	has_access_to_all_spaces = "true"
	account_role             = "Admin"
	permission               = "MANAGE_ACCOUNT"
	permissions_length       = "17"
	non_existent_user        = "nonexistent@example.com"
	timezone                 = ""
)

func TestUserDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Can't create account level tag with possible values
				Config: providerConfig + fmt.Sprintf(`
				data "torque_user" "user" {
					user_email            ="%s"
				}
				`, user_email),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.torque_user.user", "first_name", first_name),
					resource.TestCheckResourceAttr("data.torque_user.user", "last_name", last_name),
					resource.TestCheckResourceAttr("data.torque_user.user", "display_first_name", display_first_name),
					resource.TestCheckResourceAttr("data.torque_user.user", "display_last_name", display_last_name),
					resource.TestCheckResourceAttr("data.torque_user.user", "user_type", user_type),
					resource.TestCheckResourceAttr("data.torque_user.user", "join_date", join_date),
					resource.TestCheckResourceAttr("data.torque_user.user", "account_role", account_role),
					resource.TestCheckResourceAttr("data.torque_user.user", "has_access_to_all_spaces", has_access_to_all_spaces),
					resource.TestCheckResourceAttr("data.torque_user.user", "permissions.#", permissions_length),
					resource.TestCheckResourceAttr("data.torque_user.user", "permissions.0", permission),
					resource.TestCheckResourceAttr("data.torque_user.user", "timezone", timezone),
				),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
				data "torque_user" "user" {
					user_email = "%s"
				}
				`, non_existent_user),
				ExpectError: regexp.MustCompile(`Unable to Read Torque user`),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
