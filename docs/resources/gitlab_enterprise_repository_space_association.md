---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "torque_gitlab_enterprise_repository_space_association Resource - terraform-provider-torque"
subcategory: ""
description: |-
  Onboard a new GitlabEnterprise repository into an existing space
---

# torque_gitlab_enterprise_repository_space_association (Resource)

Onboard a new GitlabEnterprise repository into an existing space

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

resource "torque_gitlab_enterprise_repository_space_association" "repository" {
  space_name      = "space_name"
  base_url        = "base_url"
  repository_name = "repository_name"
  repository_url  = "repository_url"
  token           = "token"
  branch          = "branch"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `base_url` (String) Repository base URL. For example: https://gitlab-on-prem.example.com/
- `branch` (String) Repository branch to use for blueprints and automation assets
- `repository_name` (String) The name of the GitlabEnterprise repository to onboard. In this example, repo_name
- `repository_url` (String) The url of the specific GitlabEnterprise repository/project to onboard. For example: https://gitlab-on-prem.example.com/repo_name
- `space_name` (String) Existing Torque Space name
- `token` (String) Authentication Token to the project/repository