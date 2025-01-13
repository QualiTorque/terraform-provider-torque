// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestSpacesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Can't create account level tag with possible values
				Config: providerConfig + `
				data "torque_spaces" "torque_spaces" {
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.torque_spaces.torque_spaces",
						tfjsonpath.New("spaces"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.torque_spaces.torque_spaces",
						tfjsonpath.New("spaces").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.torque_spaces.torque_spaces",
						tfjsonpath.New("spaces").AtSliceIndex(0).AtMapKey("icon"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.torque_spaces.torque_spaces",
						tfjsonpath.New("spaces").AtSliceIndex(0).AtMapKey("color"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.torque_spaces.torque_spaces",
						tfjsonpath.New("spaces").AtSliceIndex(0).AtMapKey("num_of_users"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.torque_spaces.torque_spaces",
						tfjsonpath.New("spaces").AtSliceIndex(0).AtMapKey("num_of_groups"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}
