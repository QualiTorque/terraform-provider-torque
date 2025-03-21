---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "torque_servicenow_approval_channel Resource - terraform-provider-torque"
subcategory: ""
description: |-
  Creates a new ServiceNow approval channel.
---

# torque_servicenow_approval_channel (Resource)

Creates a new ServiceNow approval channel.

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

resource "torque_servicenow_approval_channel" "channel" {
  name        = "name"
  description = "description"
  base_url    = "base_url"
  username    = "username"
  password    = "password"
  headers     = "json"
  approver    = "approver@company.com"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `approver` (String) ServiceNow Approver
- `base_url` (String) ServiceNow Instance Base URL
- `name` (String) Name of the approval channel.
- `password` (String, Sensitive) ServiceNow Password
- `user_name` (String) ServiceNow Username

### Optional

- `description` (String) Description of the approval channel
- `headers` (String) Custom Headers (JSON) - JSON formatted string that represents the custom headers, for example {header:'val'}
