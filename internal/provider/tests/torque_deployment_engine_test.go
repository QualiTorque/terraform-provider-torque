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
		deployment_engine_name           = "argo_deployment_engine"
		description                      = "description"
		agent_name                       = "demo-prod"
		auth_token                       = "token"
		server_url                       = "https://argocd.com"
		polling_interval_seconds         = "30"
		polling_interval_seconds_int     = 30
		all_spaces                       = "true"
		new_deployment_engine_name       = "new_argo_deployment_engine"
		new_description                  = "new_description"
		new_auth_token                   = "new_token"
		new_server_url                   = "https://new_argocd.com"
		new_polling_interval_seconds     = "60"
		specific_spaces                  = "TorqueTerraformProvider"
		new_polling_interval_seconds_int = 60
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
					statecheck.ExpectKnownValue(
						"torque_deployment_engine.engine",
						tfjsonpath.New("description"),
						knownvalue.StringExact(description),
					),
					statecheck.ExpectKnownValue(
						"torque_deployment_engine.engine",
						tfjsonpath.New("agent_name"),
						knownvalue.StringExact(agent_name),
					),
					statecheck.ExpectKnownValue(
						"torque_deployment_engine.engine",
						tfjsonpath.New("server_url"),
						knownvalue.StringExact(server_url),
					),
					statecheck.ExpectKnownValue(
						"torque_deployment_engine.engine",
						tfjsonpath.New("auth_token"),
						knownvalue.StringExact(auth_token),
					),
					statecheck.ExpectKnownValue(
						"torque_deployment_engine.engine",
						tfjsonpath.New("polling_interval_seconds"),
						knownvalue.Int32Exact(polling_interval_seconds_int),
					),
				},
			},
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_deployment_engine" "engine" {
					name                     = "%s"
					description              = "%s"
					agent_name               = "%s"
					auth_token               = "%s"
					server_url               = "%s"
					polling_interval_seconds = "%s"
					specific_spaces          = ["%s"]
					// all_spaces               = "%s"
				}
				`, new_deployment_engine_name, new_description, agent_name, new_auth_token, new_server_url, new_polling_interval_seconds, specific_spaces, all_spaces),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"torque_deployment_engine.engine",
						tfjsonpath.New("name"),
						knownvalue.StringExact(new_deployment_engine_name),
					),
					statecheck.ExpectKnownValue(
						"torque_deployment_engine.engine",
						tfjsonpath.New("description"),
						knownvalue.StringExact(new_description),
					),
					statecheck.ExpectKnownValue(
						"torque_deployment_engine.engine",
						tfjsonpath.New("agent_name"),
						knownvalue.StringExact(agent_name),
					),
					statecheck.ExpectKnownValue(
						"torque_deployment_engine.engine",
						tfjsonpath.New("server_url"),
						knownvalue.StringExact(new_server_url),
					),
					statecheck.ExpectKnownValue(
						"torque_deployment_engine.engine",
						tfjsonpath.New("auth_token"),
						knownvalue.StringExact(new_auth_token),
					),
					statecheck.ExpectKnownValue(
						"torque_deployment_engine.engine",
						tfjsonpath.New("polling_interval_seconds"),
						knownvalue.Int32Exact(new_polling_interval_seconds_int),
					),
				},
			},
		},
	})
}
