---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "torque_workflow Resource - terraform-provider-torque"
subcategory: ""
description: |-
  Allows to enable and publish existing Torque workflow with env or env_resource scope so it will be allowed to be executed and displayed in the self-service catalog.
---

# torque_workflow (Resource)

Allows to enable and publish existing Torque workflow with env or env_resource scope so it will be allowed to be executed and displayed in the self-service catalog.

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


resource "torque_workflow" "env_workflow" {
  name            = "Day2"
  space_name      = "Space"
  repository_name = "Repo"
  self_service    = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the workflow to enable.
- `repository_name` (String) Repository where the workflow source code is
- `space_name` (String) Space the workflow belongs to

### Optional

- `self_service` (Boolean) Indicates whether this workflow is displayed in the self-service catalog. For workflows with Space scope, then this field can be omitted and will always be true.

### Read-Only

- `launch_allowed` (Boolean) Indicates whether this workflow is enabled and allowed to be launched