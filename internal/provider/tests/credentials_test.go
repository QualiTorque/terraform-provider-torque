// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccountCredentials(t *testing.T) {
	const (
		name                  = "credentials"
		description           = "description"
		token                 = "token"
		git_type              = "github"
		space                 = "TorqueTerraformProvider"
		new_name              = "new_credentials_new"
		new_description       = "new_description"
		new_token             = "new_token"
		new_git_type          = "bitbucket"
		new_space             = "TorqueTerraformProvider-2"
		allowed_spaces_length = "1"
	)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_git_credentials" "credentials" {
					name                = "%s"
					description         = "%s"
					token               = "%s"
					type                = "%s"
					allowed_space_names = ["%s"]
				}
				`, name, description, token, git_type, space),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_git_credentials.credentials", "name", name),
					resource.TestCheckResourceAttr("torque_git_credentials.credentials", "description", description),
					resource.TestCheckResourceAttr("torque_git_credentials.credentials", "token", token),
					resource.TestCheckResourceAttr("torque_git_credentials.credentials", "type", git_type),
					resource.TestCheckResourceAttr("torque_git_credentials.credentials", "allowed_space_names.0", space),
					resource.TestCheckResourceAttr("torque_git_credentials.credentials", "allowed_space_names.#", allowed_spaces_length),
				),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_git_credentials" "credentials" {
					name                = "%s"
					description         = "%s"
					token               = "%s"
					type                = "%s"
					allowed_space_names = ["%s"]
				}
				`, new_name, new_description, new_token, new_git_type, new_space),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_git_credentials.credentials", "name", new_name),
					resource.TestCheckResourceAttr("torque_git_credentials.credentials", "description", new_description),
					resource.TestCheckResourceAttr("torque_git_credentials.credentials", "token", new_token),
					resource.TestCheckResourceAttr("torque_git_credentials.credentials", "type", new_git_type),
					resource.TestCheckResourceAttr("torque_git_credentials.credentials", "allowed_space_names.0", new_space),
					resource.TestCheckResourceAttr("torque_git_credentials.credentials", "allowed_space_names.#", allowed_spaces_length),
				),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_git_credentials" "credentials" {
					name                = "%s"
					description         = "%s"
					token               = "%s"
					type                = "%s"
					// allowed_space_names = ["%s"]
				}
				`, new_name, new_description, new_token, new_git_type, new_space),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_git_credentials.credentials", "name", new_name),
					resource.TestCheckResourceAttr("torque_git_credentials.credentials", "description", new_description),
					resource.TestCheckResourceAttr("torque_git_credentials.credentials", "token", new_token),
					resource.TestCheckResourceAttr("torque_git_credentials.credentials", "type", new_git_type),
					resource.TestCheckNoResourceAttr("torque_git_credentials.credentials", "allowed_space_names"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestSpaceCredentials(t *testing.T) {
	const (
		name                  = "credentials"
		description           = "description"
		token                 = "token"
		git_type              = "github"
		space                 = "TorqueTerraformProvider"
		new_name              = "new_credentials_new"
		new_description       = "new_description"
		new_token             = "new_token"
		new_git_type          = "bitbucket"
		new_space             = "TorqueTerraformProvider-2"
		allowed_spaces_length = "1"
	)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_space_git_credentials" "credentials" {
					name                = "%s"
					description         = "%s"
					token               = "%s"
					type                = "%s"
					space_name          = "%s"
				}
				`, name, description, token, git_type, space),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_space_git_credentials.credentials", "name", name),
					resource.TestCheckResourceAttr("torque_space_git_credentials.credentials", "description", description),
					resource.TestCheckResourceAttr("torque_space_git_credentials.credentials", "token", token),
					resource.TestCheckResourceAttr("torque_space_git_credentials.credentials", "type", git_type),
					resource.TestCheckResourceAttr("torque_space_git_credentials.credentials", "space_name", space),
				),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_space_git_credentials" "credentials" {
					name                = "%s"
					description         = "%s"
					token               = "%s"
					type                = "%s"
					space_name          = "%s"
				}
				`, new_name, new_description, new_token, new_git_type, new_space),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_space_git_credentials.credentials", "name", new_name),
					resource.TestCheckResourceAttr("torque_space_git_credentials.credentials", "description", new_description),
					resource.TestCheckResourceAttr("torque_space_git_credentials.credentials", "token", new_token),
					resource.TestCheckResourceAttr("torque_space_git_credentials.credentials", "type", new_git_type),
					resource.TestCheckResourceAttr("torque_space_git_credentials.credentials", "space_name", new_space),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
