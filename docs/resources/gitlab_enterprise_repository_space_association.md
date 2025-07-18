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
  repository_name = "repository_name"
  repository_url  = "repository_url"
  token           = "token"
  branch          = "branch"
  credential_name = "credentials"
  use_all_agents  = false
  agents          = ["eks", "aks"]
  timeout         = 5
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `branch` (String) Repository branch to use for blueprints and automation assets
- `credential_name` (String) The name of the Credentials to use/create. Must be unique in the space.
- `repository_name` (String) The name of the GitlabEnterprise repository to onboard. In this example, repo_name
- `repository_url` (String) The url of the specific GitlabEnterprise repository/project to onboard. For example: https://gitlab-on-prem.example.com/repo_name
- `space_name` (String) Existing Torque Space name

### Optional

- `agents` (List of String) List of specific agents to use to onboard and sync this repository. Cannot be specified when use_all_agents is true.
- `auto_register_eac` (Boolean) Auto register environment files
- `timeout` (Number) Time in minutes to wait for Torque to sync the repository during the onboarding. Default is 1 minute.
- `token` (String, Deprecated) Authentication Token to the project/repository. If omitted, existing credentials provided in the credential_name field will be used for authentication. If provided, a new credentials object will be created.
- `use_all_agents` (Boolean) Whether all associated agents can be used to onboard and sync this repository. Must be set to false if agents attribute is used.
