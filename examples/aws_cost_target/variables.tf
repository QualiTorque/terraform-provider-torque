variable "torque_space" {
  type        = string
  description = "Torque space to be used during Torque API authentication."
}

variable "torque_token" {
  type        = string
  sensitive = true
  description = "Torque token (long token or short token) to be used for Torque API authneticaiton."
}

variable "cost_target_name" {
  type        = string
  description = "The name of the new cost target"
}

variable "role_arn" {
  type        = string
  description = "The AWS role ARN to be used used in the cost target"
}

variable "external_id" {
  type        = string
  description = "The AWS external id to be used in the cost target"
}
