module "blueprint" {
  source          = "./blueprint"
  space_name      = "space"
  blueprint_name  = "blueprint"
  repository_name = "repository_name"
  self_service    = true
  building_block  = true
  labels          = ["aws", "k8s"]
  tags = {
    "activity_type" = "demo"
  }
  always_on        = true
  allow_scheduling = false
  default_icon     = "nodes"
  custom_icon      = "icon"
}
