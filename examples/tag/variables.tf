variable "torque_space" {
  type        = string
  description = "Torque space to be used during Torque API authentication."
}

variable "torque_token" {
  type        = string
  sensitive   = true
  description = "Torque token (long token or short token) to be used for Torque API authneticaiton."
}

variable "tag_name" {
  type        = string
  description = "The name of the tag to be created"
}

variable "tag_value" {
  type        = string
  description = "The value of the new tag"
}

variable "tag_scope" {
  type        = string
  description = "The scope of the new tag. Possible values: account, space, blueprint, environment"
}

variable "tag_description" {
  type        = string
  description = "The description to be added to the tag"
}


variable "tag_possible_values" {
  type        = list(string)
  description = "list of possible values for the tag. The tag value should be part of the possible values list."
}