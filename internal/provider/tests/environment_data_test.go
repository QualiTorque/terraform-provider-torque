// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestEnvironmentDataSource(t *testing.T) {
	const (
		id                        = "cMX1RkaWj6gm"
		name                      = "S3-TF-Provider-Test"
		blueprint_name            = "S3 Bucket"
		blueprint_commit          = "790fc72a3454f7dd1b7667669340ef157ee5a43c"
		blueprint_repository_name = "TerraformProviderAcceptanceTests"
		status                    = "Active"
		collaborators_length      = "1"
		collaborator              = "terraformtester@quali.com"
		is_eac                    = "false"
		last_used                 = "2024-08-29T08:41:24.7433562Z"
		end_time                  = ""
		start_time                = "2024-08-29T08:40:11.6257438Z"
		owner_email               = "amir.r@quali.com"
		initiator_email           = "amir.r@quali.com"
		raw_json                  = "{\"owner\":{\"first_name\":\"Amir\",\"last_name\":\"Rashkovsky\",\"timezone\":\"Asia/Jerusalem\",\"email\":\"amir.r@quali.com\",\"join_date\":\"2024-08-20T09:24:28.886853\",\"display_first_name\":\"Amir\",\"display_last_name\":\"Rashkovsky\"},\"initiator\":{\"first_name\":\"Amir\",\"last_name\":\"Rashkovsky\",\"timezone\":\"Asia/Jerusalem\",\"email\":\"amir.r@quali.com\",\"join_date\":\"2024-08-20T09:24:28.886853\",\"display_first_name\":\"Amir\",\"display_last_name\":\"Rashkovsky\"},\"collaborators_info\":{\"collaborators\":[{\"first_name\":\"Developer\",\"last_name\":\"McDeveloper\",\"timezone\":null,\"email\":\"terraformtester@quali.com\",\"join_date\":\"2024-08-23T11:20:56.827262\",\"display_first_name\":\"Developer\",\"display_last_name\":\"McDeveloper\"}],\"all_space_members\":false},\"is_workflow\":false,\"details\":{\"state\":{\"grains\":[{\"state\":{\"storage_name\":\"storage-12a3fb7\",\"stages\":[{\"id\":\"76aafe7d-5a33-48f5-abc5-88031a2f24db\",\"name\":\"Deploy\",\"execution\":{\"start\":\"2024-08-29T08:40:12.844781+00:00\",\"duration\":\"00:01:11.0009029\"},\"activities\":[{\"id\":\"75577262-8e0c-4177-b275-4a620b43a8c3\",\"name\":\"Plan Validation\",\"status\":\"Skipped\",\"log\":\"/api/environment/logs/75577262-8e0c-4177-b275-4a620b43a8c3\",\"execution\":{\"start\":null,\"duration\":null},\"errors\":[]},{\"id\":\"a65b028f-fea5-4959-8503-435eceb0a712\",\"name\":\"Prepare\",\"status\":\"Done\",\"log\":\"/api/environment/logs/a65b028f-fea5-4959-8503-435eceb0a712\",\"execution\":{\"start\":\"2024-08-29T08:40:12.8447583+00:00\",\"duration\":\"00:00:00.5045040\"},\"errors\":[]},{\"id\":\"8e0189cd-e408-41f4-9ccd-5f20ab686b13\",\"name\":\"Init\",\"status\":\"Done\",\"log\":\"/api/environment/logs/8e0189cd-e408-41f4-9ccd-5f20ab686b13\",\"execution\":{\"start\":\"2024-08-29T08:40:40.2655756+00:00\",\"duration\":\"00:00:07.9276848\"},\"errors\":[]},{\"id\":\"05f7528a-1194-4369-ad32-e2baac7dbfb2\",\"name\":\"Tagging\",\"status\":\"Done\",\"log\":\"/api/environment/logs/05f7528a-1194-4369-ad32-e2baac7dbfb2\",\"execution\":{\"start\":\"2024-08-29T08:40:48.1943742+00:00\",\"duration\":\"00:00:10.1314549\"},\"errors\":[]},{\"id\":\"45c958d7-e011-4a81-bbf4-aa83e74a585e\",\"name\":\"Plan\",\"status\":\"Done\",\"log\":\"/api/environment/logs/45c958d7-e011-4a81-bbf4-aa83e74a585e\",\"execution\":{\"start\":\"2024-08-29T08:40:58.3262441+00:00\",\"duration\":\"00:00:13.2529309\"},\"errors\":[]},{\"id\":\"5e1ad4cd-2d75-4f09-ac08-d8e624dc187c\",\"name\":\"Apply\",\"status\":\"Done\",\"log\":\"/api/environment/logs/5e1ad4cd-2d75-4f09-ac08-d8e624dc187c\",\"execution\":{\"start\":\"2024-08-29T08:41:11.6001356+00:00\",\"duration\":\"00:00:12.2114653\"},\"errors\":[]}],\"errors\":[],\"details\":null}],\"current_state\":\"Deployed\",\"drift\":{\"deployment\":{\"detected\":false},\"asset\":{\"detected\":false,\"dismissed\":false,\"deployed_commit_sha\":\"\",\"latest_commit_sha\":\"\"}}},\"id\":\"c48e1970-69d6-4fa4-933e-90d8673a459b\",\"name\":\"s3\",\"path\":\"s3\",\"kind\":\"terraform\",\"imported\":false,\"execution_host\":\"demo-prod\",\"details\":{\"backend\":null,\"type\":\"terraform\"},\"sources\":[{\"store\":\"TerraformProviderAcceptanceTests\",\"path\":\"assets/s3\",\"full_path\":\"https://github.com/QualiNext/TerraformProviderAcceptanceTests/tree/master/assets/s3\",\"branch\":\"master\",\"commit\":\"790fc72a3454f7dd1b7667669340ef157ee5a43c\",\"commit_date\":\"08/29/2024 08:40:25\",\"author\":\"Amir Rashkovsky\",\"commit_message\":\"rename bp\",\"is_default_branch\":true,\"is_last_commit\":true}],\"workspace_directories\":[],\"depends_on\":[],\"inputs\":[]}],\"current_state\":\"active\",\"execution\":{\"retention\":{\"kind\":\"indefinite\"},\"start_time\":\"2024-08-29T08:40:11.6257438Z\",\"end_time\":null},\"errors\":[],\"outputs\":[],\"eac_synced\":false},\"id\":\"cMX1RkaWj6gm\",\"definition\":{\"metadata\":{\"name\":\"S3-TF-Provider-Test\",\"space_name\":\"TorqueTerraformProvider\",\"automation\":false,\"eac_url\":null,\"blueprint\":\"S3 Bucket\",\"blueprint_name\":\"S3 Bucket\",\"blueprint_display_name\":\"S3 Bucket\",\"blueprint_inputs\":[{\"name\":\"agent\",\"type\":\"agent\",\"style\":\"text\",\"default_value\":null,\"has_default_value\":false,\"sensitive\":false,\"description\":null,\"allowed_values\":[],\"parameter_name\":null,\"pattern\":null,\"validation_description\":null,\"depends_on\":[],\"source_name\":null,\"overrides\":[]}],\"blueprint_commit\":\"790fc72a3454f7dd1b7667669340ef157ee5a43c\",\"repository_name\":\"TerraformProviderAcceptanceTests\"},\"inputs\":[{\"name\":\"agent\",\"value\":\"demo-prod\"}],\"instructions\":{\"text\":null,\"url\":null},\"layout\":null,\"tags\":[{\"name\":\"activity_type\",\"value\":\"other\",\"modified_by\":\"Amir Rashkovsky\",\"last_modified\":\"2024-08-29T08:40:10.9770352Z\",\"created_by\":\"Amir Rashkovsky\",\"created_date\":\"2024-08-29T08:40:10.9770352Z\",\"scope\":\"runtime\",\"possible_values\":[\"development\",\"manual-testing\",\"automation-testing\",\"security-testing\",\"load-testing\",\"production\",\"staging\",\"demo\",\"recreation\",\"other\"],\"description\":\"The business activity type this sandbox was launched for\",\"tag_type\":\"pre_defined\"},{\"name\":\"rewrwerew\",\"value\":\"unassigned_tag_value\",\"modified_by\":\"Yarin Keren\",\"last_modified\":\"2024-08-26T10:37:51.573432\",\"created_by\":\"Yarin Keren\",\"created_date\":\"2024-08-26T10:37:51.573432\",\"scope\":\"space\",\"possible_values\":[],\"description\":\"\",\"tag_type\":\"user_defined\"},{\"name\":\"torque-sandbox-name\",\"value\":\"S3 Bucket-20240829T11404363\",\"modified_by\":\"Torque\",\"last_modified\":null,\"created_by\":\"Torque\",\"created_date\":null,\"scope\":\"runtime\",\"possible_values\":null,\"description\":null,\"tag_type\":\"system\"},{\"name\":\"torque-space-name\",\"value\":\"TorqueTerraformProvider\",\"modified_by\":\"Torque\",\"last_modified\":null,\"created_by\":\"Torque\",\"created_date\":null,\"scope\":\"runtime\",\"possible_values\":null,\"description\":null,\"tag_type\":\"system\"},{\"name\":\"torque-blueprint-name\",\"value\":\"TerraformProviderAcceptanceTests/S3 Bucket\",\"modified_by\":\"Torque\",\"last_modified\":null,\"created_by\":\"Torque\",\"created_date\":null,\"scope\":\"runtime\",\"possible_values\":null,\"description\":null,\"tag_type\":\"system\"},{\"name\":\"torque-owner-email\",\"value\":\"amir.r@quali.com\",\"modified_by\":\"Torque\",\"last_modified\":null,\"created_by\":\"Torque\",\"created_date\":null,\"scope\":\"runtime\",\"possible_values\":null,\"description\":null,\"tag_type\":\"system\"},{\"name\":\"torque-environment-id\",\"value\":\"cMX1RkaWj6gm\",\"modified_by\":\"Torque\",\"last_modified\":null,\"created_by\":\"Torque\",\"created_date\":null,\"scope\":\"runtime\",\"possible_values\":null,\"description\":null,\"tag_type\":\"system\"},{\"name\":\"torque-env-case-ignored-id\",\"value\":\"cpgiock1ppoiq4tsc7so\",\"modified_by\":\"Torque\",\"last_modified\":null,\"created_by\":\"Torque\",\"created_date\":null,\"scope\":\"runtime\",\"possible_values\":null,\"description\":null,\"tag_type\":\"system\"},{\"name\":\"torque-account-id\",\"value\":\"857ef9bb-16e9-4cbc-80af-de28c240cef5\",\"modified_by\":\"Torque\",\"last_modified\":null,\"created_by\":\"Torque\",\"created_date\":null,\"scope\":\"runtime\",\"possible_values\":null,\"description\":null,\"tag_type\":\"system\"}],\"labels\":[]},\"computed_status\":\"Active\",\"estimated_launch_duration_in_seconds\":null},\"cost\":null,\"read_only\":false,\"last_used\":\"2024-08-29T08:41:24.7433562Z\",\"annotations\":[],\"entity_metadata\":null,\"workflow_instantiation_name\":null}"
		grains_length             = "1"
		inputs_length             = "1"
		outputs_length            = "0"
		tag_name                  = "activity_type"
		tags_length               = "9"
		errors_length             = "0"
		input_name                = "agent"
		grain_kind                = "terraform"
		grain_id                  = "c48e1970-69d6-4fa4-933e-90d8673a459b"
		grain_path                = "s3"
		grain_state               = "Deployed"
		grain_name                = "s3"
		grain_source_store        = repository_name
		grain_source_path         = "assets/s3"
		grain_branch              = "master"
		grain_commit              = "790fc72a3454f7dd1b7667669340ef157ee5a43c"
		grain_is_last_commit      = "true"
	)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create and Read testing
				Config: providerConfig + fmt.Sprintf(`
					data "torque_environment" "env" {
						space_name = "%s"
						id         = "%s"
					}
				`, space_name, id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.torque_environment.env", "name", name),
					resource.TestCheckResourceAttr("data.torque_environment.env", "id", id),
					resource.TestCheckResourceAttr("data.torque_environment.env", "blueprint_name", blueprint_name),
					resource.TestCheckResourceAttr("data.torque_environment.env", "blueprint_commit", blueprint_commit),
					resource.TestCheckResourceAttr("data.torque_environment.env", "blueprint_repository_name", blueprint_repository_name),
					resource.TestCheckResourceAttr("data.torque_environment.env", "status", status),
					resource.TestCheckResourceAttr("data.torque_environment.env", "collaborators.#", collaborators_length),
					resource.TestCheckResourceAttr("data.torque_environment.env", "collaborators.0.email", collaborator),
					resource.TestCheckResourceAttr("data.torque_environment.env", "is_eac", is_eac),
					resource.TestCheckResourceAttr("data.torque_environment.env", "last_used", last_used),
					resource.TestCheckResourceAttr("data.torque_environment.env", "end_time", end_time),
					resource.TestCheckResourceAttr("data.torque_environment.env", "start_time", start_time),
					resource.TestCheckResourceAttr("data.torque_environment.env", "owner_email", owner_email),
					resource.TestCheckResourceAttr("data.torque_environment.env", "initiator_email", initiator_email),
					resource.TestCheckResourceAttr("data.torque_environment.env", "raw_json", raw_json),

					resource.TestCheckResourceAttr("data.torque_environment.env", "grains.#", grains_length),
					resource.TestCheckResourceAttr("data.torque_environment.env", "grains.0.name", grain_name),
					resource.TestCheckResourceAttr("data.torque_environment.env", "grains.0.kind", grain_kind),
					resource.TestCheckResourceAttr("data.torque_environment.env", "grains.0.id", grain_id),
					resource.TestCheckResourceAttr("data.torque_environment.env", "grains.0.path", grain_path),
					resource.TestCheckResourceAttr("data.torque_environment.env", "grains.0.state.current_state", grain_state),

					resource.TestCheckResourceAttr("data.torque_environment.env", "grains.0.sources.0.store", repository_name),
					resource.TestCheckResourceAttr("data.torque_environment.env", "grains.0.sources.0.path", grain_source_path),
					resource.TestCheckResourceAttr("data.torque_environment.env", "grains.0.sources.0.branch", grain_branch),
					resource.TestCheckResourceAttr("data.torque_environment.env", "grains.0.sources.0.commit", grain_commit),
					resource.TestCheckResourceAttr("data.torque_environment.env", "grains.0.sources.0.is_last_commit", grain_is_last_commit),

					resource.TestCheckResourceAttr("data.torque_environment.env", "inputs.#", inputs_length),
					resource.TestCheckResourceAttr("data.torque_environment.env", "inputs.0.name", input_name),
					resource.TestCheckResourceAttr("data.torque_environment.env", "outputs.#", outputs_length),
					resource.TestCheckResourceAttr("data.torque_environment.env", "tags.#", tags_length),
					resource.TestCheckResourceAttr("data.torque_environment.env", "tags.0.name", tag_name),
					resource.TestCheckResourceAttr("data.torque_environment.env", "errors.#", errors_length),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
