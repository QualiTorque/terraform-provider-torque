variable "torque_space" {
  type        = string
  description = "Torque space to be used during Torque API authentication."
}

variable "torque_token" {
  type        = string
  sensitive   = true
  description = "Torque token (long token or short token) to be used for Torque API authneticaiton."
}

variable "space_name" {
  type        = string
  description = "Torque space where the new tag should set on."
}

variable "tag_name" {
  type        = string
  description = "The tag name to be used"
}

variable "tag_value" {
  type        = string
  description = "The tag value to be set in the space"
}
