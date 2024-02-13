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
  description = "Torque space name to be created"
}

variable "space_color" {
  type        = string
  default     = "darkGreen"
  description = "The color of the new space"
}

variable "space_icon" {
  type        = string
  default     = "flow"
  description = "The icno of the new space"
}

