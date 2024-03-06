terraform {
  required_providers {
    torque = {
      source = "qualitorque/torque"
    }
  }
}

provider "torque" {
  host  = "https://portal.qtorque.io/"
  space = var.torque_space
  token = var.torque_token
}

resource "torque_space_email_notification" "new_notification" {
  space_name                    = var.space_name
  notification_name             = var.notification_name
  environment_launched          = var.environment_launched
  environment_deployed          = var.environment_deployed
  environment_force_ended       = var.environment_force_ended
  environment_idle              = var.environment_idle
  environment_extended          = var.environment_extended
  drift_detected                = var.drift_detected
  workflow_failed               = var.workflow_failed
  workflow_started              = var.workflow_started
  updates_detected              = var.updates_detected
  collaborator_added            = var.collaborator_added
  action_failed                 = var.action_failed
  environment_ending_failed     = var.environment_ending_failed
  environment_ended             = var.environment_ended
  environment_active_with_error = var.environment_active_with_error
}