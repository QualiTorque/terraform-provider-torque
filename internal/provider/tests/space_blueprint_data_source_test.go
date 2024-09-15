// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSpaceBlueprintResource(t *testing.T) {
	const (
		blueprint_test_name          = "blueprint_data_source_test"
		blueprint_display_name       = "Data Source Acceptance Test Blueprint"
		repository_name              = "TerraformProviderAcceptanceTests"
		repository_branch            = "master"
		commit                       = "7cbf9f3b10eced6b9abe0fc8a71aca697076db5d"
		blueprint_description        = "Blueprint Data Source Test"
		url                          = "https://github.com/QualiNext/TerraformProviderAcceptanceTests/blob/master/blueprints/blueprint_data_source_test.yaml"
		modified_by                  = "amiros89"
		last_modified                = "2024-08-29T04:51:19Z"
		enabled                      = "false"
		always_on                    = "false"
		default_duration             = "P8DT2H2M"
		num_of_active_environments   = "0"
		default_extend               = "PT5H1M"
		max_duration                 = "P10DT2H15M"
		max_active_environments      = "2"
		tag_name                     = "activity_type"
		default_value                = "other"
		possible_values_length       = "10"
		tag_description              = "The business activity type this sandbox was launched for"
		input_name                   = "region"
		input_default_value          = "eu-west-1"
		input_possible_values_length = "0"
		input_description            = "The AWS region where this instance will be launched"
		inputs_length                = "2"
	)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create and Read testing
				Config: providerConfig + fmt.Sprintf(`
				data "torque_space_blueprint" "blueprint" {
					name       = "%s"
					space_name = "%s"
				}
				`, blueprint_test_name, space_name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.torque_space_blueprint.blueprint", "name", blueprint_test_name),
					resource.TestCheckResourceAttr("data.torque_space_blueprint.blueprint", "space_name", space_name),
					resource.TestCheckResourceAttr("data.torque_space_blueprint.blueprint", "display_name", blueprint_display_name),
					resource.TestCheckResourceAttr("data.torque_space_blueprint.blueprint", "description", blueprint_description),
					resource.TestCheckResourceAttr("data.torque_space_blueprint.blueprint", "repository_name", repository_name),
					resource.TestCheckResourceAttr("data.torque_space_blueprint.blueprint", "repository_branch", repository_branch),
					resource.TestCheckResourceAttr("data.torque_space_blueprint.blueprint", "commit", commit),
					resource.TestCheckResourceAttr("data.torque_space_blueprint.blueprint", "url", url),
					resource.TestCheckResourceAttr("data.torque_space_blueprint.blueprint", "modified_by", modified_by),
					resource.TestCheckResourceAttr("data.torque_space_blueprint.blueprint", "last_modified", last_modified),
					resource.TestCheckResourceAttr("data.torque_space_blueprint.blueprint", "num_of_active_environments", num_of_active_environments),
					resource.TestCheckResourceAttr("data.torque_space_blueprint.blueprint", "max_duration", max_duration),
					resource.TestCheckResourceAttr("data.torque_space_blueprint.blueprint", "default_duration", default_duration),
					resource.TestCheckResourceAttr("data.torque_space_blueprint.blueprint", "default_extend", default_extend),
					resource.TestCheckResourceAttr("data.torque_space_blueprint.blueprint", "max_active_environments", max_active_environments),
					resource.TestCheckResourceAttr("data.torque_space_blueprint.blueprint", "always_on", always_on),
					resource.TestCheckResourceAttr("data.torque_space_blueprint.blueprint", "tags.#", "1"),
					resource.TestCheckResourceAttr("data.torque_space_blueprint.blueprint", "tags.0.name", tag_name),
					resource.TestCheckResourceAttr("data.torque_space_blueprint.blueprint", "tags.0.default_value", default_value),
					resource.TestCheckResourceAttr("data.torque_space_blueprint.blueprint", "tags.0.possible_values.#", possible_values_length),
					resource.TestCheckResourceAttr("data.torque_space_blueprint.blueprint", "tags.0.description", tag_description),
					resource.TestCheckResourceAttr("data.torque_space_blueprint.blueprint", "inputs.#", inputs_length),
					resource.TestCheckResourceAttr("data.torque_space_blueprint.blueprint", "inputs.0.name", input_name),
					resource.TestCheckResourceAttr("data.torque_space_blueprint.blueprint", "inputs.0.default_value", input_default_value),
					resource.TestCheckResourceAttr("data.torque_space_blueprint.blueprint", "inputs.0.possible_values.#", input_possible_values_length),
					resource.TestCheckResourceAttr("data.torque_space_blueprint.blueprint", "inputs.0.description", input_description),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
