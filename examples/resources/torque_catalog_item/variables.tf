variable "torque_space" {
  type        = string
  description = "Torque space to be used during Torque API authentication."
}

variable "torque_token" {
  type        = string
  sensitive   = true
  description = "Torque token (long token or short token) to be used for Torque API authneticaiton."
}

variable "torque_host" {
  type        = string
  default     = "https://portal.qtorque.io/"
  description = "Torque portal URL"
}

variable "space_name" {
  type        = string
  description = "Torque space where the catalog item should be published in."
}

variable "blueprint_name" {
  type        = string
  default     = "my-test-group"
  description = "The blueprint that shuld be publish to a catalog item"
}

variable "repository_name" {
  type        = string
  default     = "my-test-group"
  description = "The repository name where the blueprint is stored"
}

