terraform {
  required_providers {
    torque = {
      source = "qualitorque/torque"
    }
  }
}

provider "torque" {
  host  = "https://portal.qtorque.io/"
  space = "MuSpace"
  token = var.torque_token
}

resource "torque_space_resource" "new_space" {
  space_name = "newspacetest"
  color      = "darkGreen"
  icon       = "flow"
}

resource "torque_agent_space_association" "agenta" {
  space_name      = torque_space_resource.new_space.space_name
  namespace       = "vido-sb"
  service_account = "vido-sb"
  agent_name      = "review2-aks"
}

resource "torque_user_space_association" "users" {
  // support multiple users with the same role.
  //for_each   = toset(var.space_members)
  user_email = "quali@quali.com"
  space_name = torque_space_resource.new_space.space_name
  user_role  = "Space Member"
}

resource "torque_repository_space_association" "github" {
  // depands-on space
  space_name      = torque_space_resource.new_space.space_name
  repository_url  = "https://github.com/QualiTorque/Torque-Samples"
  access_token    = var.repo_token
  repository_type = "github"
  branch          = "main"
  repository_name = "Torque-Samples"
}

resource "torque_tag" "newtag" {
  name            = "my new tag"
  value           = "tag value"
  scope           = "space"
  possible_values = ["a", "b", "c"]
  description     = "something to say about this tag"
}

resource "torque_space_tag_value_association" "dev_space" {
  space_name  = "development"
  tag_name    = "cost_cetner"
  value_value = "Dev"
}

data "torque_space_blueprints" "blueprints" {
  space_name = "01-Development"
  //filter by repo
}

resource "torque_catalog_item_resource" "something" {
  // foreach and filter if needed
  space_name      = "01-Development"
  blueprint_name  = "eks-cluster_1"
  repository_name = "qtorque"
}

data "torque_space_tags" "tags" {
  space_name = "something"
}

response will be 
[
  { 
    "name":
    "value"
    "possible_values"
    "scope"
    "description"
  },

]