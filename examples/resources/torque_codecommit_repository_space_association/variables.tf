variable "space" {
  type        = string
  description = "Torque space to be used during Torque API authentication."
}

variable "token" {
  type        = string
  sensitive   = true
  description = "Torque token (long token or short token) to be used for Torque API authneticaiton."
}

variable "host" {
  type        = string
  default     = "https://portal.qtorque.io/"
  description = "Torque portal URL"
}

variable "aws_region" {
  type        = string
  description = "AWS region"
}

variable "role_arn" {
  type        = string
  description = "AWS Role ARN"
}

variable "external_id" {
  type        = string
  description = "External ID"
}

variable "username" {
  type        = string
  description = "Username"
}

variable "password" {
  type        = string
  description = "Password"
}

variable "branch" {
  type        = string
  description = "Branch"
}

variable "repository_name" {
  type        = string
  description = "Name of the repository"
}

variable "repository_url" {
  type        = string
  description = "URL of the repository"
}
