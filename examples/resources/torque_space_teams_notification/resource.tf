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


resource "torque_space_teams_notification" "new_notification" {
  space_name                    = "space"
  notification_name             = "notification_name"
  webhook                       = "https://webhook.com"
  environment_launched          = false
  environment_deployed          = false
  environment_force_ended       = false
  environment_idle              = false
  environment_extended          = false
  drift_detected                = false
  workflow_failed               = true
  workflow_started              = true
  updates_detected              = true
  collaborator_added            = true
  action_failed                 = false
  environment_ending_failed     = true
  environment_ended             = true
  environment_active_with_error = true
  blueprint_published           = true
  blueprint_unpublished         = true
  idle_reminders                = [1, 2, 3]
}
