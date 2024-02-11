variable "torque_space" {
  type    = string
}

variable "torque_token" {
  type    = string
  sensitive = true
}

variable "group_name" {
  type    = string
  default = "my-test-group"
}


variable "group_description" {
  type    = string
  default = ""
}

variable "group_idp_identifier" {
  type    = string
  default = ""
}

variable "group_account_role" {
  type    = string
  default = ""
}

variable "group_users" {
  type    = list(string)
  default = []
}

variable "group_users" {
  type    = list(string)
  default = []
}

variable "group_space_roles" {
  type = list(object({
    space_name = string
    space_role = string
  }))
  default = []
}