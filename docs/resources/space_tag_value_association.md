---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "torque_space_tag_value_association Resource - terraform-provider-torque"
subcategory: ""
description: |-
  Associate Torque space with existing registered agent
---

# torque_space_tag_value_association (Resource)

Associate Torque space with existing registered agent



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `space_name` (String) Existing Torque Space name
- `tag_name` (String) Tag name configured in the account
- `tag_value` (String) The tag value to be set for the space
