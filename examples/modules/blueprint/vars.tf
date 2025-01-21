variable "tags" {
  type    = map(string)
  default = null
}

variable "labels" {
  type    = list(string)
  default = []
}

variable "building_block" {
  type    = bool
  default = false
}

variable "self_service" {
  type    = bool
  default = false
}

variable "always_on" {
  type    = bool
  default = false
}

variable "allow_scheduling" {
  type    = bool
  default = false
}

variable "max_active_environments" {
  type    = number
  default = null
}

variable "max_duration" {
  type    = string
  default = "PT2H"
}

variable "default_duration" {
  type    = string
  default = "PT2H"
}

variable "default_extend" {
  type    = string
  default = "PT2H"
}

variable "default_icon" {
  type    = string
  default = "nodes"
}

variable "custom_icon" {
  type    = string
  default = null
}

variable "space_name" {
  type = string
}

variable "blueprint_name" {
  type = string
}

variable "display_name" {
  type    = string
  default = null
}

variable "repository_name" {
  type    = string
  default = "qtorque"
}
