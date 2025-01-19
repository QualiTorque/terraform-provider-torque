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

func TestTorqueDeploymentEngineResource(t *testing.T) {
	const (
		deployment_engine_name   = "argo_deployment_engine"
		description              = "description"
		agent_name               = "demo-prod"
		auth_token               = "token"
		server_url               = "https://argocd.com"
		polling_interval_seconds = "30"
		all_spaces               = "true"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_deployment_engine" "engine" {
					name                     = "%s"
					description              = "%s"
					agent_name               = "%s"
					auth_token               = "%s"
					server_url               = "%s"
					polling_interval_seconds = "%s"
					all_spaces               = "%s"
				}

				`, deployment_engine_name, description, agent_name, auth_token, server_url, polling_interval_seconds, all_spaces),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"torque_deployment_engine.engine",
						tfjsonpath.New("name"),
						knownvalue.StringExact(deployment_engine_name),
					),
				},
			},
			// Update and Read testing
			// {
			// 	Config: providerConfig + fmt.Sprintf(`
			// 	resource "torque_space_label" "test" {
			// 		space_name = "%s"
			// 		name       = "%s"
			// 		color      = "bordeaux"
			// 		quick_filter = "true"
			// 	}
			// 	`, space_name, newLabelName),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("torque_space_label.test", "space_name", space_name),
			// 		resource.TestCheckResourceAttr("torque_space_label.test", "name", newLabelName),
			// 		resource.TestCheckResourceAttr("torque_space_label.test", "color", "bordeaux"),
			// 		resource.TestCheckResourceAttr("torque_space_label.test", "quick_filter", "true"),
			// 	),
			// },
			// Delete testing automatically occurs in TestCase
		},
	})
}
