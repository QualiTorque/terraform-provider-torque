// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	workflow_name          = "TerraformWorkflowTest"
	workflow_display_name  = "Test Terraform Provider Workflow Data Source"
	workflow_description   = "It's just a test..."
	workflow_enforced      = "false"
	yaml                   = "on:\n  overridable: true\n  scheduler:\n  - 13 1 * * 3,5\ninputs: {}\njobs:\n  TerraformWorkflowTest:\n    name: Test Terraform Provider Workflow Data Source\n    steps:\n    - name: terminate-environment\n      uses: terminate-environment\n"
	specific_spaces_length = "1"
)

func TestWorkflowDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create and Read testing
				Config: providerConfig + fmt.Sprintf(`
				data "torque_workflow" "workflow" {
					name            = "%s"
				}
				`, workflow_name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.torque_workflow.workflow", "name", workflow_name),
					resource.TestCheckResourceAttr("data.torque_workflow.workflow", "yaml", yaml),
					resource.TestCheckResourceAttr("data.torque_workflow.workflow", "display_name", workflow_display_name),
					resource.TestCheckResourceAttr("data.torque_workflow.workflow", "description", workflow_description),
					resource.TestCheckResourceAttr("data.torque_workflow.workflow", "enforced_on_all_spaces", workflow_enforced),
					resource.TestCheckResourceAttr("data.torque_workflow.workflow", "specific_spaces.#", specific_spaces_length),
					resource.TestCheckResourceAttr("data.torque_workflow.workflow", "specific_spaces.0", space_name),
				),
			},
		},
	})
}
