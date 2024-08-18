// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestIntrospectionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
				resource "torque_introspection_resource" "test" {
					display_name       = "My Resource"
					image              = "https://cdn-icons-png.flaticon.com/512/882/882730.png"
					introspection_data = {size = "large", mode = "party"}
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_introspection_resource.test", "display_name", "My Resource"),
					resource.TestCheckResourceAttr("torque_introspection_resource.test", "image", "https://cdn-icons-png.flaticon.com/512/882/882730.png"),
					resource.TestCheckResourceAttr("torque_introspection_resource.test", "introspection_data.size", "large"),
					resource.TestCheckResourceAttr("torque_introspection_resource.test", "introspection_data.mode", "party"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + `
				resource "torque_introspection_resource" "test" {
					display_name       = "Another Display Name"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_introspection_resource.test", "display_name", "Another Display Name"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
