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

func TestTorqueEmailApprovalChannel(t *testing.T) {
	const (
		approval_channel = "email_approval_channel"
		description      = "description"
		new_description  = "new_description"
		approver         = "terraformtester@quali.com"
		approver2        = "terraformtester2@quali.com"
	)

	var unique_name = approval_channel + "_" + index
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_email_approval_channel" "channel" {
					name                     = "%s"
					description              = "%s"
					approvers               = ["%s"]
				}
				`, unique_name, description, approver),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"torque_email_approval_channel.channel",
						tfjsonpath.New("name"),
						knownvalue.StringExact(unique_name),
					),
					statecheck.ExpectKnownValue(
						"torque_email_approval_channel.channel",
						tfjsonpath.New("description"),
						knownvalue.StringExact(description),
					),
					statecheck.ExpectKnownValue(
						"torque_email_approval_channel.channel",
						tfjsonpath.New("approvers").AtSliceIndex(0),
						knownvalue.StringExact(approver),
					),
				},
			},
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_email_approval_channel" "channel" {
					name                     = "%s"
					description              = "%s"
					approvers               = ["%s","%s"]
				}
				`, unique_name, new_description, approver, approver2),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"torque_email_approval_channel.channel",
						tfjsonpath.New("name"),
						knownvalue.StringExact(unique_name),
					),
					statecheck.ExpectKnownValue(
						"torque_email_approval_channel.channel",
						tfjsonpath.New("description"),
						knownvalue.StringExact(new_description),
					),
					statecheck.ExpectKnownValue(
						"torque_email_approval_channel.channel",
						tfjsonpath.New("approvers").AtSliceIndex(0),
						knownvalue.StringExact(approver),
					),
					statecheck.ExpectKnownValue(
						"torque_email_approval_channel.channel",
						tfjsonpath.New("approvers").AtSliceIndex(1),
						knownvalue.StringExact(approver2),
					),
				},
			},
		},
	})
}
