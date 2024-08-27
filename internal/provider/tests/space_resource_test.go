// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSpaceResource(t *testing.T) {
	randomSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	spaceName := fmt.Sprintf("MySpace-%s", randomSuffix)
	newSpaceName := fmt.Sprintf("MyNewSpace-%s", randomSuffix)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_space" "test" {
					space_name = "%s"
					icon       = "re"
					color      = "darkBlue"
				}
				`, spaceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_space.test", "space_name", spaceName),
					resource.TestCheckResourceAttr("torque_space.test", "icon", "re"),
					resource.TestCheckResourceAttr("torque_space.test", "color", "darkBlue"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_space" "test" {
					space_name = "%s"
					icon       = "star"
					color      = "pinkRed"
				}
				`, newSpaceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_space.test", "space_name", newSpaceName),
					resource.TestCheckResourceAttr("torque_space.test", "icon", "star"),
					resource.TestCheckResourceAttr("torque_space.test", "color", "pinkRed"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
