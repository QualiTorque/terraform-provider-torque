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


resource "torque_environment" "example" {
  blueprint_name   = "Hello World"
  environment_name = "Sample Environment"
  owner_email      = "someone@company.com"
  duration         = "PT2H" # ISO 8601 duration format - must not be specified together with scheduled_end_time. Both can be omitted to create always-on environment.
  space            = "MySpace"
  force_destroy    = false
  inputs = {
    "agent" = "playground",
    "name"  = "amir"
  }
  description = "This is a sample environment created using Torque."
  tags = {
    "activity_type" = "demo"
  }
  collaborators = {
    collaborators_emails = []
    all_space_members    = true
  }
  blueprint_source = {
    blueprint_name  = "custom-blueprint"
    repository_name = "custom-repo"
    branch          = "main"
    commit          = "abcd1234"
  }
  scheduled_end_time = "2024-12-31T23:59:59Z" # must not be specified together with duration. Both can be omitted to create always-on environment.
  workflows = [                               # built-in Torque Workflows
    {
      name = "sample-workflow"

      schedules = [
        {
          scheduler  = "0 12 * * *" # CRON expression for scheduling
          overridden = true
        }
      ]

      reminder = 15 # in minutes before execution

      inputs_overrides = {
        "key1" = "value1"
        "key2" = "value2"
      }
    }
  ]
}
