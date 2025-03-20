// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestEnvironmentLabelAssociationResource(t *testing.T) {
	i, _ := strconv.Atoi(index)
	environmentIDs := []string{
		"v0b0DgneI8wt",
		"N3DTp1UPkW2z",
		"QC6DaQ3vIX7S",
		"wVGhj0g1g03u",
		"gWGHfvmHnMQg",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_environment_label_association" "test" {
					space_name = "%s"
					environment_id       = "%s"
					labels = [{"key":"test","value":"test"}]
				}
				`, space_name, environmentIDs[i]),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"torque_environment_label_association.test",
						tfjsonpath.New("space_name"),
						knownvalue.StringExact(space_name),
					),
					statecheck.ExpectKnownValue(
						"torque_environment_label_association.test",
						tfjsonpath.New("environment_id"),
						knownvalue.StringExact(environmentIDs[i]),
					),
					statecheck.ExpectKnownValue(
						"torque_environment_label_association.test",
						tfjsonpath.New("labels").AtSliceIndex(0).AtMapKey("key"),
						knownvalue.StringExact("test"),
					),
					statecheck.ExpectKnownValue(
						"torque_environment_label_association.test",
						tfjsonpath.New("labels").AtSliceIndex(0).AtMapKey("value"),
						knownvalue.StringExact("test"),
					),
					statecheck.ExpectKnownValue(
						"torque_environment_label_association.test",
						tfjsonpath.New("labels"),
						knownvalue.ListSizeExact(1),
					),
				},
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_environment_label_association" "test" {
					space_name = "%s"
					environment_id       = "%s"
					labels = [{"key":"key","value":"val"}]

				}
				`, space_name, environmentIDs[i]),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"torque_environment_label_association.test",
						tfjsonpath.New("space_name"),
						knownvalue.StringExact(space_name),
					),
					statecheck.ExpectKnownValue(
						"torque_environment_label_association.test",
						tfjsonpath.New("environment_id"),
						knownvalue.StringExact(environmentIDs[i]),
					),
					statecheck.ExpectKnownValue(
						"torque_environment_label_association.test",
						tfjsonpath.New("labels").AtSliceIndex(0).AtMapKey("key"),
						knownvalue.StringExact("key"),
					),
					statecheck.ExpectKnownValue(
						"torque_environment_label_association.test",
						tfjsonpath.New("labels").AtSliceIndex(0).AtMapKey("value"),
						knownvalue.StringExact("val"),
					),
					statecheck.ExpectKnownValue(
						"torque_environment_label_association.test",
						tfjsonpath.New("labels"),
						knownvalue.ListSizeExact(1),
					),
				},
			},
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_environment_label_association" "test" {
					space_name = "%s"
					environment_id       = "%s"
					labels = [{"key":"key","value":"val"},{"key":"test","value":"test"}]
				}
				`, space_name, environmentIDs[i]),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"torque_environment_label_association.test",
						tfjsonpath.New("space_name"),
						knownvalue.StringExact(space_name),
					),
					statecheck.ExpectKnownValue(
						"torque_environment_label_association.test",
						tfjsonpath.New("environment_id"),
						knownvalue.StringExact(environmentIDs[i]),
					),
					statecheck.ExpectKnownValue(
						"torque_environment_label_association.test",
						tfjsonpath.New("labels").AtSliceIndex(0).AtMapKey("key"),
						knownvalue.StringExact("key"),
					),
					statecheck.ExpectKnownValue(
						"torque_environment_label_association.test",
						tfjsonpath.New("labels").AtSliceIndex(0).AtMapKey("value"),
						knownvalue.StringExact("val"),
					),
					statecheck.ExpectKnownValue(
						"torque_environment_label_association.test",
						tfjsonpath.New("labels").AtSliceIndex(1).AtMapKey("key"),
						knownvalue.StringExact("test"),
					),
					statecheck.ExpectKnownValue(
						"torque_environment_label_association.test",
						tfjsonpath.New("labels").AtSliceIndex(1).AtMapKey("value"),
						knownvalue.StringExact("test"),
					),
					statecheck.ExpectKnownValue(
						"torque_environment_label_association.test",
						tfjsonpath.New("labels"),
						knownvalue.ListSizeExact(2),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
