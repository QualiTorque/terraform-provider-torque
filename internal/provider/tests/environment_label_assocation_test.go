// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestEnvironmentLabelAssociationResource(t *testing.T) {
	const (
		environment_id = "HLjN9yHyCYeD"
	)

	// randomSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	// label := fmt.Sprintf("EnvLabel-%s", randomSuffix)

	// newLabelName := fmt.Sprintf("MyNewEnvLabel-%s", randomSuffix)
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
				`, space_name, environment_id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_environment_label_association.test", "space_name", space_name),
					resource.TestCheckResourceAttr("torque_environment_label_association.test", "environment_id", environment_id),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_environment_label_association" "test" {
					space_name = "%s"
					environment_id       = "%s"
					labels = [{"key":"test1","value":"test1"}]
				}
				`, space_name, environment_id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_environment_label_association.test", "space_name", space_name),
					resource.TestCheckResourceAttr("torque_environment_label_association.test", "environment_id", environment_id),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
