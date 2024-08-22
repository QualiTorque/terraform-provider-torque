// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/qualitorque/terraform-provider-torque/client"
)

const (
	blueprint_name     = "ec2"
	repository_name    = "TerraformProviderAcceptanceTests"
	new_blueprint_name = "rds"
)

// minorVresion       = (strings.Split(os.Getenv("INDEX"), "."))[1]

var version = os.Getenv("VERSION")
var minorVresion = strings.Split((version), ".")
var index = minorVresion[1]
var unique_blueprint_name = blueprint_name + "_" + index
var new_unique_blueprint_name = new_blueprint_name + "_" + index

func TestCatalogItemResource(t *testing.T) {
	spaceName := os.Getenv("TORQUE_SPACE")
	fmt.Println(unique_blueprint_name)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testBlueprintNotPublished(new_unique_blueprint_name),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_catalog_item" "catalog_item" {
					space_name      = "%s"
					blueprint_name  = "%s"
					repository_name = "%s"
				}
				`, spaceName, unique_blueprint_name, repository_name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_catalog_item.catalog_item", "space_name", spaceName),
					resource.TestCheckResourceAttr("torque_catalog_item.catalog_item", "blueprint_name", unique_blueprint_name),
					resource.TestCheckResourceAttr("torque_catalog_item.catalog_item", "repository_name", repository_name),
					testBlueprintPublished(unique_blueprint_name),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_catalog_item" "catalog_item" {
					space_name      = "%s"
					blueprint_name  = "%s"
					repository_name = "%s"
				}
				`, spaceName, new_unique_blueprint_name, repository_name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_catalog_item.catalog_item", "space_name", spaceName),
					resource.TestCheckResourceAttr("torque_catalog_item.catalog_item", "blueprint_name", new_unique_blueprint_name),
					resource.TestCheckResourceAttr("torque_catalog_item.catalog_item", "repository_name", repository_name),
					testBlueprintPublished(new_unique_blueprint_name),
				),
				// 	ExpectError: regexp.MustCompile("Unable to publish blueprint in space"),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestCatalogItemErrorIfNotExists(t *testing.T) {
	spaceName := os.Getenv("TORQUE_SPACE")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_catalog_item" "catalog_item" {
					space_name      = "%s"
					blueprint_name  = "non-existent-blueprint"
					repository_name = "non-existent-repo"
				}
				`, spaceName),
				ExpectError: regexp.MustCompile("Unable to publish blueprint in space"),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testBlueprintPublished(blueprint string) resource.TestCheckFunc {
	return checkBlueprintPublishedCondition(true, blueprint)
}

func testBlueprintNotPublished(blueprint string) resource.TestCheckFunc {
	return checkBlueprintPublishedCondition(false, blueprint)
}

func checkBlueprintPublishedCondition(expectedPublished bool, blueprint string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		host := os.Getenv("TORQUE_HOST")
		space := os.Getenv("TORQUE_SPACE")
		token := os.Getenv("TORQUE_TOKEN")
		c, err := client.NewClient(&host, &space, &token)
		if err != nil {
			return err
		}
		bp, err := c.GetBlueprint(space, blueprint)
		if err != nil {
			return err
		}
		if bp.Published != expectedPublished {
			return fmt.Errorf("expected Published to be '%v', got '%v'", expectedPublished, bp.Published)
		}

		return nil
	}
}
