---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "torque_agent_space_association Resource - terraform-provider-torque"
subcategory: ""
description: |-
  Associate Torque space with existing registered agent
---

# torque_agent_space_association (Resource)

Associate Torque space with existing registered agent

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

resource "torque_agent_space_association" "agent_association" {
  space_name      = "space"
  agent_name      = "agent"
  service_account = "service-account"
  namespace       = "default"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `agent_name` (String) Agent name to associate to the space
- `namespace` (String) Default namespace to be used with the agent in the space
- `service_account` (String) Default service account to be used with the agent in the space
- `space_name` (String) Existing Torque Space name
