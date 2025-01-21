// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestTorqueServiceNowApprovalChannel(t *testing.T) {
	const (
		approval_channel = "approval_channel"
		description      = "description"
		new_description  = "new_description"
		approver         = "terraformtester@quali.com"
		approver2        = "terraformtester2@quali.com"
		base_url         = "base_url"
		new_base_url     = "new_base_url"
		user_name        = "username"
		password         = "password"
		new_user_name    = "new_username"
		new_password     = "new_password"
	)

	var unique_name = approval_channel + "_" + index
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_servicenow_approval_channel" "channel" {
					name                    = "%s"
					description             = "%s"
					base_url                = "%s"
					user_name 				= "%s"
					password                = "%s"
					approver                = "%s"
				}
				`, unique_name, description, base_url, user_name, password, approver),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"torque_servicenow_approval_channel.channel",
						tfjsonpath.New("name"),
						knownvalue.StringExact(unique_name),
					),
					statecheck.ExpectKnownValue(
						"torque_servicenow_approval_channel.channel",
						tfjsonpath.New("description"),
						knownvalue.StringExact(description),
					),
					statecheck.ExpectKnownValue(
						"torque_servicenow_approval_channel.channel",
						tfjsonpath.New("base_url"),
						knownvalue.StringExact(base_url),
					),
					statecheck.ExpectKnownValue(
						"torque_servicenow_approval_channel.channel",
						tfjsonpath.New("approver"),
						knownvalue.StringExact(approver),
					),
				},
			},
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_servicenow_approval_channel" "channel" {
					name                    = "%s"
					description             = "%s"
					base_url                = "%s"
					user_name 				= "%s"
					password                = "%s"
					approver                = "%s"
				}
				`, unique_name, new_description, new_base_url, new_user_name, new_password, approver2),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"torque_servicenow_approval_channel.channel",
						tfjsonpath.New("name"),
						knownvalue.StringExact(unique_name),
					),
					statecheck.ExpectKnownValue(
						"torque_servicenow_approval_channel.channel",
						tfjsonpath.New("description"),
						knownvalue.StringExact(new_description),
					),
					statecheck.ExpectKnownValue(
						"torque_servicenow_approval_channel.channel",
						tfjsonpath.New("base_url"),
						knownvalue.StringExact(new_base_url),
					),
					statecheck.ExpectKnownValue(
						"torque_servicenow_approval_channel.channel",
						tfjsonpath.New("approver"),
						knownvalue.StringExact(approver2),
					),
				},
			},
			// {
			// 	Config: providerConfig + fmt.Sprintf(`
			// 	resource "torque_servicenow_approval_channel" "channel" {
			// 		name                    = "%s"
			// 		description             = "%s"
			// 		base_url                = "%s"
			// 		user_name = 				"%s"
			// 		password =              "%s"
			// 		approvers               = ["%s","%s"]
			// 	}
			// 	`, unique_name, description, base_url, username, password, approver, approver2),
			// 	ConfigStateChecks: []statecheck.StateCheck{
			// 		statecheck.ExpectKnownValue(
			// 			"torque_teams_approval_channel.channel",
			// 			tfjsonpath.New("name"),
			// 			knownvalue.StringExact(unique_name),
			// 		),
			// 		statecheck.ExpectKnownValue(
			// 			"torque_teams_approval_channel.channel",
			// 			tfjsonpath.New("description"),
			// 			knownvalue.StringExact(new_description),
			// 		),
			// 		statecheck.ExpectKnownValue(
			// 			"torque_teams_approval_channel.channel",
			// 			tfjsonpath.New("webhook_address"),
			// 			knownvalue.StringExact(new_webhook_address),
			// 		),
			// 		statecheck.ExpectKnownValue(
			// 			"torque_teams_approval_channel.channel",
			// 			tfjsonpath.New("approvers").AtSliceIndex(0),
			// 			knownvalue.StringExact(approver),
			// 		),
			// 		statecheck.ExpectKnownValue(
			// 			"torque_teams_approval_channel.channel",
			// 			tfjsonpath.New("approvers").AtSliceIndex(1),
			// 			knownvalue.StringExact(approver2),
			// 		),
			// 	},
			// },
		},
	})
}
