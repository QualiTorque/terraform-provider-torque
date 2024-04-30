terraform {
  required_providers {
    torque = {
      source = "qualitorque/torque"
    }
  }
}

provider "torque" {
  host  = "https://portal.qtorque.io/"
  space = "space"
  token = "111111111111"
}
resource "torque_environment" "name" {
  space          = "space"
  blueprint_name = "blueprint"
  duration       = "PT2H"
  automation     = false
  description    = "value"
  inputs = {
    "name" = "value"
  }
  environment_name = "hello"
  owner_email      = "owner@email.com"
  tags = {
    "activity_type" = "development"
  }
  workflows = [{
    name = "workflow"
    schedules = [{
      scheduler  = " 5 0 * 8 *"
      overridden = false
    }]
    inputs_overrides = {
      "input1" = "value1"
    }
    reminder = 3
  }]

  collaborators = {
    collaborators_emails = ["collaborator@email.com"]
    all_space_members    = false
  }

  blueprint_source = {
    blueprint_name  = "blueprint"
    repository_name = "my repository"
    branch          = "master"
    commit          = "commit"
  }

}
