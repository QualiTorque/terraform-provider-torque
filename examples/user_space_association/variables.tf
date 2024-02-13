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
  description = "Torque space where the user will be added to"
}

variable "user_email" {
  type        = string
  description = "User email that is already signed up for Torque"
}

variable "user_role" {
  type        = string
  description = "The desire role of the user in the space"
}
