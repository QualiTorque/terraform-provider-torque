---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "torque_tag Resource - terraform-provider-torque"
subcategory: ""
description: |-
  Creation of a new Torque space with associated entities (users, repos, etc...)
---

# torque_tag (Resource)

Creation of a new Torque space with associated entities (users, repos, etc...)



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the new tag to be added to torque
- `scope` (String) Tag scope. Possible values: account, space, blueprint, environment
- `value` (String) Tag value to be set as the tag value default

### Optional

- `description` (String) Tag description
- `possible_values` (List of String) Tag possible values
