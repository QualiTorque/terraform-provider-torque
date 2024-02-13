variable "torque_space" {
  type        = string
  description = "Torque space to be used during Torque API authentication."
}

variable "torque_token" {
  type        = string
  sensitive = true
  description = "Torque token (long token or short token) to be used for Torque API authneticaiton."
}

variable "parameter_name" {
  type        = string
  description = "The name of the account paraemter to be created."
}

variable "parameter_value" {
  type        = string
  sensitive   = true
  description = "The value of the account paramter to be used"
}

variable "parameter_sensitive" {
  type        = bool
  default     = false
  description = "Set the account parameter as sensitive so it's value will be hidden from users"
}


variable "parameter_description" {
  type        = string
  description = "The parameter description to be presented in the torque user interface"
}
