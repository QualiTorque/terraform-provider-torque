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

func TestSpaceTeamsNotificationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_space_teams_notification" "notification" {
					space_name                    = "%s"
					notification_name             = "%s"
					web_hook                      = "%s"
					environment_launched          = false
					environment_deployed          = false
					environment_force_ended       = false
					environment_idle              = false
					environment_extended          = false
					drift_detected                = false
					workflow_failed               = true
					workflow_started              = true
					updates_detected              = true
					collaborator_added            = true
					action_failed                 = false
					environment_ending_failed     = true
					environment_ended             = true
					environment_active_with_error = true
					blueprint_published           = true
					blueprint_unpublished         = true
					idle_reminders                = [1, 2, 3]
				}
				`, fullSpaceName, notificationName, web_hook),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("space_name"),
						knownvalue.StringExact(fullSpaceName),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("notification_name"),
						knownvalue.StringExact(notificationName),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("web_hook"),
						knownvalue.StringExact(web_hook),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("environment_launched"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("environment_deployed"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("environment_force_ended"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("environment_idle"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("environment_extended"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("drift_detected"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("workflow_failed"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("workflow_started"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("updates_detected"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("collaborator_added"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("action_failed"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("environment_ending_failed"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("environment_ended"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("environment_active_with_error"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("blueprint_published"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("blueprint_unpublished"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("idle_reminders"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.Int32Exact(1),
							knownvalue.Int32Exact(2),
							knownvalue.Int32Exact(3),
						}),
					),
				},
			},
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_space_teams_notification" "notification" {
					space_name                    = "%s"
					notification_name             = "%s"
					web_hook                      = "%s"
					environment_launched          = true
					environment_deployed          = true
					environment_force_ended       = true
					environment_idle              = true
					environment_extended          = true
					drift_detected                = true
					workflow_failed               = false
					workflow_started              = false
					updates_detected              = false
					collaborator_added            = false
					action_failed                 = false
					environment_ending_failed     = false
					environment_ended             = false
					environment_active_with_error = false
					blueprint_published           = false
					blueprint_unpublished         = true
					idle_reminders                = [1, 2]
				}
				`, fullSpaceName, notificationName, web_hook),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("space_name"),
						knownvalue.StringExact(fullSpaceName),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("notification_name"),
						knownvalue.StringExact(notificationName),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("web_hook"),
						knownvalue.StringExact(web_hook),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("environment_launched"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("environment_deployed"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("environment_force_ended"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("environment_idle"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("environment_extended"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("drift_detected"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("workflow_failed"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("workflow_started"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("updates_detected"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("collaborator_added"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("action_failed"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("environment_ending_failed"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("environment_ended"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("environment_active_with_error"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("blueprint_published"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("blueprint_unpublished"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("idle_reminders"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.Int32Exact(1),
							knownvalue.Int32Exact(2),
						}),
					),
					statecheck.ExpectKnownValue(
						"torque_space_teams_notification.notification",
						tfjsonpath.New("idle_reminders"),
						knownvalue.ListSizeExact(2),
					),
				},
			},
		},
	})
}
