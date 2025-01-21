resource "torque_catalog_item" "name" {
  count           = var.self_service ? 1 : 0
  space_name      = var.space_name
  blueprint_name  = var.blueprint_name
  repository_name = var.repository_name
}

resource "torque_asset_library_item" "name" {
  count           = var.building_block ? 1 : 0
  space_name      = var.space_name
  blueprint_name  = var.blueprint_name
  repository_name = var.repository_name
}

resource "torque_space_label_association" "label_association" {
  space_name      = var.space_name
  blueprint_name  = var.blueprint_name
  repository_name = var.repository_name
  labels          = var.labels
}

resource "torque_blueprint_tag_value_association" "tag_association" {
  for_each       = var.tags
  space_name     = var.space_name
  repo_name      = var.repository_name
  tag_name       = each.key
  tag_value      = each.value
  blueprint_name = var.blueprint_name
}

resource "torque_blueprint_consumption_policy" "consumption_policy" {
  blueprint_name          = var.blueprint_name
  space_name              = var.space_name
  repo_name               = var.repository_name
  max_duration            = var.max_duration
  default_duration        = var.default_duration
  default_extend          = var.default_extend
  max_active_environments = var.max_active_environments # 10
  always_on               = var.always_on               # false
  allow_scheduling        = var.allow_scheduling        # true
}

resource "torque_blueprint_display_name" "display_name" {
  count          = var.display_name != null ? 1 : 0
  space_name     = var.space_name
  repo_name      = var.repository_name
  blueprint_name = var.blueprint_name
  display_name   = var.display_name
}

resource "torque_blueprint_default_icon" "default_icon" {
  space_name     = var.space_name
  repo_name      = var.repository_name
  blueprint_name = var.blueprint_name
  icon           = var.default_icon
}

resource "torque_blueprint_custom_icon" "custom_icon" {
  count          = var.custom_icon != null ? 1 : 0
  space_name     = var.space_name
  repo_name      = var.repository_name
  blueprint_name = var.blueprint_name
  icon           = var.custom_icon
}

# resource "torque_blueprint_display_settings" "name" {
#   space_name     = var.space_name
#   repo_name      = var.repository_name
#   blueprint_name = var.blueprint_name
#   display_name   = "display_name"
#   default_icon   = "icon"
#   custom_icon    = "custom_icon"
# }
