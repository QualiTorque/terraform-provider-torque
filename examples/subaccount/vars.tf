variable "torque_space" {
  type        = string
  description = "Torque space to be used during Torque API authentication."
}

variable "torque_token" {
  type        = string
  sensitive   = true
  description = "Torque token (long token or short token) to be used for Torque API authneticaiton."
}

variable "parent_account" {
  type        = string
  description = "Name of the Torque parent account, which is the main account of the sub-account to be added to Torque."
}

variable "account_name" {
  type        = string
  description = "Name of the new Torque sub-account to be added to Torque."
}

variable "password" {
  type        = string
  description = "Torque sub-account password"
  sensitive   = true
}

variable "company" {
  type        = string
  description = "Company name which the sub-account belongs to."
}
