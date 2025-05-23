---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "torque_elasticsearch_audit Resource - terraform-provider-torque"
subcategory: ""
description: |-
  Creates an Elasticsearch Audit target. Once integrated, Torque begins capturing events and you’ll ship them to the configured Elasticsearch instance.
---

# torque_elasticsearch_audit (Resource)

Creates an Elasticsearch Audit target. Once integrated, Torque begins capturing events and you’ll ship them to the configured Elasticsearch instance.

## Example Usage

```terraform
terraform {
  required_providers {
    torque = {
      source = "qualitorque/torque"
    }
  }
}

# provider "torque" {
#   host  = "https://portal.qtorque.io/"
#   space = "space"
#   token = "111111111111"
# }

resource "torque_elasticsearch_audit" "audit" {
  url      = "https://elastic:9000"
  username = "user"
  password = "password"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `password` (String, Sensitive) Elasticsearch instance password.
- `url` (String) Elasticsearch instance URL.
- `username` (String) Elasticsearch instance username.

### Optional

- `certificate` (String, Sensitive) Optional certificate of the Elasticsearch instance.

### Read-Only

- `type` (String) Built-in audit is of type elasticsearch.
