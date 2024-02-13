variable "torque_space" {
  type        = string
  description = "Torque space used for authenticating into Torque"
}

variable "torque_token" {
  type        = string
  sensitive   = true
  description = "Torque token used to authenticate to the Torque API. Supports both short and long tokens"
}

variable "group_name" {
  type        = string
  default     = "my-test-group"
  description = "The new Torque group name to be created"
}


variable "group_description" {
  type        = string
  default     = ""
  description = "The new Torque group description to be used"
}

variable "group_idp_identifier" {
  type        = string
  default     = ""
  description = "The Identity Provider ID to be used with the new Torque group"
}

variable "group_account_role" {
  type        = string
  default     = ""
  description = "Choose this to set the group as an account level group. Possible values: \"Admin\" or \"Member\""
}

variable "group_users" {
  type        = list(string)
  default     = []
  description = "List of user emails to assign to the group"
}

variable "group_space_roles" {
  type = list(object({
    space_name = string
    space_role = string
  }))
  default     = []
  description = "Use this to set the group scope for specific spaces and the role of the group within the spaces."
}