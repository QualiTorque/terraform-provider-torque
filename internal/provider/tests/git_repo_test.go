// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	repository_url  = "https://github.com/QualiTorque/terraform-provider-torque"
	branch          = "master"
	credential_name = "TerraformGitCreds"
	repository_type = "github"
	repo_name       = "terraform-provider-torque"
)

func TestGitRepoWithCreds(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_repository_space_association" "repository_with_credentials" {
					space_name      = "%s"
					repository_url  = "%s"
					repository_type = "%s"
					branch          = "%s"
					repository_name = "%s"
					credential_name = "%s"
				}
				`, fullSpaceName, repository_url, repository_type, branch, repo_name, credential_name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_repository_space_association.repository_with_credentials", "space_name", fullSpaceName),
					resource.TestCheckResourceAttr("torque_repository_space_association.repository_with_credentials", "repository_name", repo_name),
					resource.TestCheckResourceAttr("torque_repository_space_association.repository_with_credentials", "branch", branch),
					resource.TestCheckResourceAttr("torque_repository_space_association.repository_with_credentials", "repository_url", repository_url),
					resource.TestCheckResourceAttr("torque_repository_space_association.repository_with_credentials", "repository_type", repository_type),
				),
			},
			// Update and Read testing
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestGitRepoWithToken(t *testing.T) {
	access_token := os.Getenv("GITHUB_TOKEN")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_repository_space_association" "repository_with_token" {
					space_name      = "%s"
					repository_url  = "%s"
					access_token    = "%s"
					repository_type = "%s"
					branch          = "%s"
					repository_name = "%s"
				}
				`, fullSpaceName, repository_url, access_token, repository_type, branch, repo_name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_repository_space_association.repository_with_token", "space_name", fullSpaceName),
					resource.TestCheckResourceAttr("torque_repository_space_association.repository_with_token", "repository_name", repo_name),
					resource.TestCheckResourceAttr("torque_repository_space_association.repository_with_token", "branch", branch),
					resource.TestCheckResourceAttr("torque_repository_space_association.repository_with_token", "repository_url", repository_url),
					resource.TestCheckResourceAttr("torque_repository_space_association.repository_with_token", "repository_type", repository_type),
				),
			},
			// Update and Read testing
			// Delete testing automatically occurs in TestCase
		},
	})
}
