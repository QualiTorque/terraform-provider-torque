---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "torque_space_email_notification Resource - terraform-provider-torque"
subcategory: ""
description: |-
  Creation of a new notification is a Torque space
---

# torque_space_email_notification (Resource)

Creation of a new notification is a Torque space

## Example Usage

```terraform
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


resource "torque_space_email_notification" "new_notification" {
  space_name                    = "space"
  notification_name             = "notification_name"
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
  idle_reminders                = [1, 2, 3]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `notification_name` (String) The notification cofngiuration name in the space
- `space_name` (String) Space name to add the notification to

### Optional

- `action_failed` (Boolean) Confgure notification for the "Action Failed" event
- `collaborator_added` (Boolean) Confgure notification for the "Collaborator Added" event
- `drift_detected` (Boolean) Confgure notification for the "Drift Detected" event
- `end_threshold` (Number) Confgure notification environment end reminder notification
- `environment_active_with_error` (Boolean) Confgure notification for the "Environment Active With Error" event
- `environment_deployed` (Boolean) Confgure notification for the "Environment Deployed" event
- `environment_ended` (Boolean) Confgure notification for the "Environment Ended" event
- `environment_ending_failed` (Boolean) Confgure notification for the "Environment Ending Failed" event
- `environment_extended` (Boolean) Confgure notification for the "Environment Extended" event
- `environment_force_ended` (Boolean) Confgure notification for the "Environment Force Ended" event
- `environment_idle` (Boolean) Confgure notification for the "Environment Idle" event
- `environment_launched` (Boolean) Confgure notification for the "Environment Launched" event
- `idle_reminders` (List of Number) Confgure notification for environment reminder notification
- `updates_detected` (Boolean) Confgure notification for the "Updates Detected" event
- `workflow_failed` (Boolean) Confgure notification for the "Workflow Failed" event
- `workflow_start_reminder` (Number) Confgure notification for the "Drift Detected" event
- `workflow_started` (Boolean) Confgure notification for the "Workflow Started" event

### Read-Only

- `notification_id` (String) The id of the newly added notification
