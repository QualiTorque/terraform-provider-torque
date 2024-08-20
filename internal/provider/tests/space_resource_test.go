// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSpaceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
				resource "torque_space" "test" {
					space_name = "MySpace"
					icon       = "re"
					color      = "darkBlue"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_space.test", "space_name", "MySpace"),
					resource.TestCheckResourceAttr("torque_space.test", "icon", "re"),
					resource.TestCheckResourceAttr("torque_space.test", "color", "darkBlue"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + `
				resource "torque_space" "test" {
					space_name = "MyNewSpace"
					icon       = "star"
					color      = "pinkRed"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_space.test", "space_name", "MyNewSpace"),
					resource.TestCheckResourceAttr("torque_space.test", "icon", "star"),
					resource.TestCheckResourceAttr("torque_space.test", "color", "pinkRed"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
