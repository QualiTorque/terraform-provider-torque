variable "torque_space" {
  type        = string
  description = "Torque space to be used during Torque API authentication."
}

variable "torque_token" {
  type        = string
  sensitive = true
  description = "Torque token (long token or short token) to be used for Torque API authneticaiton."
}

variable "space_name" {
  type        = string
  description = "The space where the blueprint is stored"
}

variable "blueprint_name" {
  type        = string
  description = "The blueprint to be used"
}

variable "blueprint_repo_name" {
  type        = string
  description = "The repository name where the blueprint is stored"
}

variable "tag_name" {
  type        = bool
  default     = false
  description = "The tag name to be used"
}


variable "tag_value" {
  type        = string
  description = "The tag value to be set on the blueprint"
}
