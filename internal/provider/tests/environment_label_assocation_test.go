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
		environment_id = "F2UXQBPCcWXY"
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
					resource.TestCheckResourceAttr("torque_environment_label_association.test", "labels.#", "1"),
					resource.TestCheckResourceAttr("torque_environment_label_association.test", "labels.0.key", "test"),
					resource.TestCheckResourceAttr("torque_environment_label_association.test", "labels.0.value", "test"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_environment_label_association" "test" {
					space_name = "%s"
					environment_id       = "%s"
					labels = [{"key":"key","value":"val"}]
				}
				`, space_name, environment_id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_environment_label_association.test", "space_name", space_name),
					resource.TestCheckResourceAttr("torque_environment_label_association.test", "environment_id", environment_id),
					resource.TestCheckResourceAttr("torque_environment_label_association.test", "labels.#", "1"),
					resource.TestCheckResourceAttr("torque_environment_label_association.test", "labels.0.key", "key"),
					resource.TestCheckResourceAttr("torque_environment_label_association.test", "labels.0.value", "val"),
				),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_environment_label_association" "test" {
					space_name = "%s"
					environment_id       = "%s"
					labels = [{"key":"key","value":"val"},{"key":"test","value":"test"}]
				}
				`, space_name, environment_id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_environment_label_association.test", "space_name", space_name),
					resource.TestCheckResourceAttr("torque_environment_label_association.test", "environment_id", environment_id),
					resource.TestCheckResourceAttr("torque_environment_label_association.test", "labels.#", "2"),
					resource.TestCheckResourceAttr("torque_environment_label_association.test", "labels.0.key", "key"),
					resource.TestCheckResourceAttr("torque_environment_label_association.test", "labels.0.value", "val"),
					resource.TestCheckResourceAttr("torque_environment_label_association.test", "labels.1.key", "test"),
					resource.TestCheckResourceAttr("torque_environment_label_association.test", "labels.1.value", "test"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
