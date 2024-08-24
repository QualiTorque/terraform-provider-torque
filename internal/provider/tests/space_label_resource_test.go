// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSpaceLabelResource(t *testing.T) {
	randomSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	label := fmt.Sprintf("MyLabel-%s", randomSuffix)
	newLabelName := fmt.Sprintf("MyNewLabel-%s", randomSuffix)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_space_label" "test" {
					space_name = "%s"
					name       = "%s"
					color      = "aws"
				}
				`, space_name, label),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_space_label.test", "space_name", space_name),
					resource.TestCheckResourceAttr("torque_space_label.test", "name", label),
					resource.TestCheckResourceAttr("torque_space_label.test", "color", "aws"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_space_label" "test" {
					space_name = "%s"
					name       = "%s"
					color      = "bordeaux"
				}
				`, space_name, newLabelName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_space_label.test", "space_name", space_name),
					resource.TestCheckResourceAttr("torque_space_label.test", "name", newLabelName),
					resource.TestCheckResourceAttr("torque_space_label.test", "color", "bordeaux"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
