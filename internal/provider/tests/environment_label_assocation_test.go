// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestEnvironmentLabelAssociationResource(t *testing.T) {
	var version = os.Getenv("VERSION")
	var minorVresion = strings.Split((version), ".")
	var index_string = minorVresion[1]
	index, _ := strconv.Atoi(index_string)
	environmentIDs := []string{
		"F2UXQBPCcWXY",
		"wN2HAEV4Bte0",
		"QC6DaQ3vIX7S",
		"v0b0DgneI8wt",
		"HLjN9yHyCYeD",
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
				`, space_name, environmentIDs[index]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_environment_label_association.test", "space_name", space_name),
					resource.TestCheckResourceAttr("torque_environment_label_association.test", "environment_id", environmentIDs[index]),
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
				`, space_name, environmentIDs[index]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_environment_label_association.test", "space_name", space_name),
					resource.TestCheckResourceAttr("torque_environment_label_association.test", "environment_id", environmentIDs[index]),
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
				`, space_name, environmentIDs[index]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_environment_label_association.test", "space_name", space_name),
					resource.TestCheckResourceAttr("torque_environment_label_association.test", "environment_id", environmentIDs[index]),
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
