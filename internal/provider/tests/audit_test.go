// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

const (
	url             = "https://elastic.com:9000"
	username        = "elastic_user"
	password        = "password"
	certificate     = "cert"
	new_url         = "https://elastic.com:9001"
	new_username    = "new_elastic_user"
	new_password    = "new_password"
	new_certificate = "new_cert"
)

func TestTorqueAuditResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				resource "torque_audit" "audit" {
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("torque_audit.audit", "type", "Torque"),
				),
			},
		},
	})
}

func TestTorqueElasticsearchAuditResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_elasticsearch_audit" "audit" {
					url      = "%s"
					username = "%s"
					password = "%s"
					certificate = "%s"
				}
				`, url, username, password, certificate),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"torque_elasticsearch_audit.audit",
						tfjsonpath.New("url"),
						knownvalue.StringExact(url),
					),
					statecheck.ExpectKnownValue(
						"torque_elasticsearch_audit.audit",
						tfjsonpath.New("username"),
						knownvalue.StringExact(username),
					),
					statecheck.ExpectKnownValue(
						"torque_elasticsearch_audit.audit",
						tfjsonpath.New("password"),
						knownvalue.StringExact(password),
					),
					statecheck.ExpectKnownValue(
						"torque_elasticsearch_audit.audit",
						tfjsonpath.New("certificate"),
						knownvalue.StringExact(certificate),
					),
				},
			},
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "torque_elasticsearch_audit" "audit" {
					url      = "%s"
					username = "%s"
					password = "%s"
					certificate = "%s"
				}
				`, new_url, new_username, new_password, new_certificate),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"torque_elasticsearch_audit.audit",
						tfjsonpath.New("url"),
						knownvalue.StringExact(new_url),
					),
					statecheck.ExpectKnownValue(
						"torque_elasticsearch_audit.audit",
						tfjsonpath.New("username"),
						knownvalue.StringExact(new_username),
					),
					statecheck.ExpectKnownValue(
						"torque_elasticsearch_audit.audit",
						tfjsonpath.New("password"),
						knownvalue.StringExact(new_password),
					),
					statecheck.ExpectKnownValue(
						"torque_elasticsearch_audit.audit",
						tfjsonpath.New("certificate"),
						knownvalue.StringExact(new_certificate),
					),
				},
			},
		},
	})
}
