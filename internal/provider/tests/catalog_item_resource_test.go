// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/qualitorque/terraform-provider-torque/client"
)

const (
	blueprint_name     = "ec2"
	repository_name    = "TerraformProviderAcceptanceTests"
	new_blueprint_name = "rds"
)

func TestCatalogItemResource(t *testing.T) {
	spaceName := os.Getenv("TORQUE_SPACE")
	var unique_blueprint_name = blueprint_name + "_" + index
	var new_unique_blueprint_name = new_blueprint_name + "_" + index

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
					always_on       = "true"
				}
				`, spaceName, unique_blueprint_name, repository_name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_catalog_item.catalog_item", "space_name", spaceName),
					resource.TestCheckResourceAttr("torque_catalog_item.catalog_item", "blueprint_name", unique_blueprint_name),
					resource.TestCheckResourceAttr("torque_catalog_item.catalog_item", "repository_name", repository_name),
					resource.TestCheckResourceAttr("torque_catalog_item.catalog_item", "always_on", "true"),
					resource.TestCheckNoResourceAttr("torque_catalog_item.catalog_item", "default_duration"),
					resource.TestCheckNoResourceAttr("torque_catalog_item.catalog_item", "default_extend"),
					resource.TestCheckNoResourceAttr("torque_catalog_item.catalog_item", "max_duration"),
					resource.TestCheckNoResourceAttr("torque_catalog_item.catalog_item", "labels"),
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
					labels          = ["k8s","aws"]				
				}
				`, spaceName, new_unique_blueprint_name, repository_name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_catalog_item.catalog_item", "space_name", spaceName),
					resource.TestCheckResourceAttr("torque_catalog_item.catalog_item", "blueprint_name", new_unique_blueprint_name),
					resource.TestCheckResourceAttr("torque_catalog_item.catalog_item", "repository_name", repository_name),
					resource.TestCheckResourceAttr("torque_catalog_item.catalog_item", "always_on", "false"),
					resource.TestCheckResourceAttr("torque_catalog_item.catalog_item", "default_duration", "PT2H"),
					resource.TestCheckResourceAttr("torque_catalog_item.catalog_item", "default_extend", "PT2H"),
					resource.TestCheckResourceAttr("torque_catalog_item.catalog_item", "max_duration", "PT2H"),
					resource.TestCheckResourceAttr("torque_catalog_item.catalog_item", "labels.#", "2"),
					testBlueprintPublished(new_unique_blueprint_name),
				),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_catalog_item" "catalog_item" {
					space_name      = "%s"
					blueprint_name  = "%s"
					repository_name = "%s"
					display_name    ="display_name"
					default_duration = "PT3H"
					default_extend = "PT9H"
					max_duration = "P1DT6H"
					labels          = ["k8s"]						
				}
				`, spaceName, new_unique_blueprint_name, repository_name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_catalog_item.catalog_item", "space_name", spaceName),
					resource.TestCheckResourceAttr("torque_catalog_item.catalog_item", "blueprint_name", new_unique_blueprint_name),
					resource.TestCheckResourceAttr("torque_catalog_item.catalog_item", "repository_name", repository_name),
					resource.TestCheckResourceAttr("torque_catalog_item.catalog_item", "always_on", "false"),
					resource.TestCheckResourceAttr("torque_catalog_item.catalog_item", "default_duration", "PT3H"),
					resource.TestCheckResourceAttr("torque_catalog_item.catalog_item", "default_extend", "PT9H"),
					resource.TestCheckResourceAttr("torque_catalog_item.catalog_item", "max_duration", "P1DT6H"),
					resource.TestCheckResourceAttr("torque_catalog_item.catalog_item", "labels.#", "1"),
					testBlueprintPublished(new_unique_blueprint_name),
				),
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
				ExpectError: regexp.MustCompile("Unable to create Catalog Item"),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestCatalogItemDurationConflicts(t *testing.T) {
	spaceName := os.Getenv("TORQUE_SPACE")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_catalog_item" "catalog_item" {
					space_name      = "%s"
					blueprint_name  = "doesnt matter"
					repository_name = "doesnt matter"
					always_on = "true"
					default_duration = "PT3H"
					default_extend = "PT9H"
					max_duration = "PT30H"					
				}
				`, spaceName),
				ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
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

		maxRetries := 3
		delay := 10 * time.Second

		for i := 0; i < maxRetries; i++ {
			bp, err := c.GetBlueprint(space, blueprint)
			if err != nil {
				return err
			}
			if bp.Published == expectedPublished {
				return nil
			}

			if i < maxRetries-1 {
				time.Sleep(delay)
			}
		}
		return fmt.Errorf("expected Published to be '%v', got '%v'", expectedPublished, !expectedPublished)
	}
}
