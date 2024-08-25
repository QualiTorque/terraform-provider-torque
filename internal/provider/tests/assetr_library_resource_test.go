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

// const (
// 	blueprint_name     = "ec2"
// 	repository_name    = "TerraformProviderAcceptanceTests"
// 	new_blueprint_name = "rds"
// )

func TestAssetLibraryItemResource(t *testing.T) {
	spaceName := os.Getenv("TORQUE_SPACE")
	var version = os.Getenv("VERSION")
	var minorVresion = strings.Split((version), ".")
	var index = minorVresion[1]
	var unique_blueprint_name = blueprint_name + "_" + index
	var new_unique_blueprint_name = new_blueprint_name + "_" + index

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testBlueprintNotInAssetLibrary(new_unique_blueprint_name),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_asset_library_item" "library_item" {
					space_name      = "%s"
					blueprint_name  = "%s"
					repository_name = "%s"
				}
				`, spaceName, unique_blueprint_name, repository_name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_asset_library_item.library_item", "space_name", spaceName),
					resource.TestCheckResourceAttr("torque_asset_library_item.library_item", "blueprint_name", unique_blueprint_name),
					resource.TestCheckResourceAttr("torque_asset_library_item.library_item", "repository_name", repository_name),
					testBlueprintInAssetLibrary(unique_blueprint_name),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_asset_library_item" "library_item" {
					space_name      = "%s"
					blueprint_name  = "%s"
					repository_name = "%s"
				}
				`, spaceName, new_unique_blueprint_name, repository_name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_asset_library_item.library_item", "space_name", spaceName),
					resource.TestCheckResourceAttr("torque_asset_library_item.library_item", "blueprint_name", new_unique_blueprint_name),
					resource.TestCheckResourceAttr("torque_asset_library_item.library_item", "repository_name", repository_name),
					testBlueprintInAssetLibrary(new_unique_blueprint_name),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})

}

func TestLibraryItemErrorIfNotExists(t *testing.T) {
	spaceName := os.Getenv("TORQUE_SPACE")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_asset_library_item" "library_item" {
					space_name      = "%s"
					blueprint_name  = "non-existent-blueprint"
					repository_name = "non-existent-repo"
				}
				`, spaceName),
				ExpectError: regexp.MustCompile("Unable to add blueprint to asset-library"),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testBlueprintInAssetLibrary(blueprint string) resource.TestCheckFunc {
	return checkBlueprintAssetLibraryCondition(true, blueprint)
}

func testBlueprintNotInAssetLibrary(blueprint string) resource.TestCheckFunc {
	return checkBlueprintAssetLibraryCondition(false, blueprint)
}

func checkBlueprintAssetLibraryCondition(expectedInAssetLibrary bool, blueprint string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		host := os.Getenv("TORQUE_HOST")
		space := os.Getenv("TORQUE_SPACE")
		token := os.Getenv("TORQUE_TOKEN")

		c, err := client.NewClient(&host, &space, &token)
		if err != nil {
			return err
		}
		bp, err := c.GetBlueprintFromAssetLibrary(space, blueprint)
		if bp != nil && expectedInAssetLibrary {
			fmt.Println("found blueprint in asset library and expected it to be there")
			return nil
		}
		if bp == nil && !expectedInAssetLibrary {
			fmt.Println("did not find blueprint in asset library and did not expect it to be there")
			return nil
		}

		return fmt.Errorf("expected blueprint in asset-library condition to be %v, got %s", expectedInAssetLibrary, err)
	}
}
