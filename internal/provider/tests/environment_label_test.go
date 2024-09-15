// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestEnvironmentLabelResource(t *testing.T) {
	const (
		label_value     = "value"
		new_label_value = "new_value"
	)

	randomSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	label := fmt.Sprintf("EnvLabel-%s", randomSuffix)

	newLabelName := fmt.Sprintf("MyNewEnvLabel-%s", randomSuffix)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_environment_label" "test" {
					key = "%s"
					value       = "%s"
				}
				`, label, label_value),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_environment_label.test", "key", label),
					resource.TestCheckResourceAttr("torque_environment_label.test", "value", label_value),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_environment_label" "test" {
					key = "%s"
					value       = "%s"
				}
				`, newLabelName, new_label_value),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_environment_label.test", "key", newLabelName),
					resource.TestCheckResourceAttr("torque_environment_label.test", "value", new_label_value),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
