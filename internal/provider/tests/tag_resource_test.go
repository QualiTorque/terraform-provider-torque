// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

// func TestTagResource(t *testing.T) {
// 	randomSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
// 	tagName := fmt.Sprintf("tag_name_%s", randomSuffix)
// 	newTagName := fmt.Sprintf("new_tag_name_%s", randomSuffix)
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { testAccPreCheck(t) },
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				// Can't create account level tag with possible values
// 				Config: providerConfig + fmt.Sprintf(`
// 				resource "torque_tag" "tag" {
// 					name            = "%s"
// 					value           = "tag_value"
// 					scope           = "account"
// 					description     = "tag_description"
// 					possible_values = ["value1", "value2"]
// 				}
// 				`, tagName),
// 				ExpectError: regexp.MustCompile("Unable to create tag"),
// 			},
// 			{
// 				// Create and Read testing
// 				Config: providerConfig + fmt.Sprintf(`
// 				resource "torque_tag" "tag" {
// 					name            = "%s"
// 					value           = "tag_value"
// 					scope           = "account"
// 					description     = "tag_description"
// 				}
// 				`, tagName),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestCheckResourceAttr("torque_tag.tag", "name", tagName),
// 					resource.TestCheckResourceAttr("torque_tag.tag", "value", "tag_value"),
// 					resource.TestCheckResourceAttr("torque_tag.tag", "scope", "account"),
// 					resource.TestCheckResourceAttr("torque_tag.tag", "description", "tag_description"),
// 					resource.TestCheckResourceAttr("torque_tag.tag", "possible_values.#", "0"),
// 				),
// 			},
// 			// Update and Read testing
// 			{
// 				Config: providerConfig + fmt.Sprintf(`
// 				resource "torque_tag" "tag" {
// 					name            = "%s"
// 					value           = "new_tag_value"
// 					scope           = "space"
// 					description     = "new_tag_description"
// 					possible_values = ["value1", "value2"]
// 				}
// 				`, newTagName),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestCheckResourceAttr("torque_tag.tag", "name", newTagName),
// 					resource.TestCheckResourceAttr("torque_tag.tag", "value", "new_tag_value"),
// 					resource.TestCheckResourceAttr("torque_tag.tag", "scope", "space"),
// 					resource.TestCheckResourceAttr("torque_tag.tag", "description", "new_tag_description"),
// 					resource.TestCheckResourceAttr("torque_tag.tag", "possible_values.0", "value1"),
// 					resource.TestCheckResourceAttr("torque_tag.tag", "possible_values.1", "value2"),
// 					resource.TestCheckResourceAttr("torque_tag.tag", "possible_values.#", "2"),
// 				),
// 			},
// 			// Delete testing automatically occurs in TestCase
// 		},
// 	})
// }
