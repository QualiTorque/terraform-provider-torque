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

const (
	agentName         = "demo-prod"
	serviceAccount    = "service_account"
	namespace         = "default"
	newServiceAccount = "new_service_account"
	newNamespace      = "new_ns"
)

func TestAgentSpaceAssocationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Can't create account level tag with possible values
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_agent_space_association" "agent_association" {
					space_name      = "%s"
					agent_name      = "%s"
					service_account = "%s"
					namespace       = "%s"
				}
				`, fullSpaceName, agentName, serviceAccount, newNamespace),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"torque_agent_space_association.agent_association",
						tfjsonpath.New("space_name"),
						knownvalue.StringExact(fullSpaceName),
					),
					statecheck.ExpectKnownValue(
						"torque_agent_space_association.agent_association",
						tfjsonpath.New("agent_name"),
						knownvalue.StringExact(agentName),
					),
					statecheck.ExpectKnownValue(
						"torque_agent_space_association.agent_association",
						tfjsonpath.New("service_account"),
						knownvalue.StringExact(serviceAccount),
					),
					statecheck.ExpectKnownValue(
						"torque_agent_space_association.agent_association",
						tfjsonpath.New("namespace"),
						knownvalue.StringExact(newNamespace),
					),
				},
			},
			{
				// Can't create account level tag with possible values
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_agent_space_association" "agent_association" {
					space_name      = "%s"
					agent_name      = "%s"
					service_account = "%s"
					namespace       = "%s"
				}
				`, fullSpaceName, agentName, newServiceAccount, newNamespace),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"torque_agent_space_association.agent_association",
						tfjsonpath.New("space_name"),
						knownvalue.StringExact(fullSpaceName),
					),
					statecheck.ExpectKnownValue(
						"torque_agent_space_association.agent_association",
						tfjsonpath.New("agent_name"),
						knownvalue.StringExact(agentName),
					),
					statecheck.ExpectKnownValue(
						"torque_agent_space_association.agent_association",
						tfjsonpath.New("service_account"),
						knownvalue.StringExact(newServiceAccount),
					),
					statecheck.ExpectKnownValue(
						"torque_agent_space_association.agent_association",
						tfjsonpath.New("namespace"),
						knownvalue.StringExact(newNamespace),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
