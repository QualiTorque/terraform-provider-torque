variable "torque_space" {
  type        = string
  description = "Torque space to be used during Torque API authentication."
}

variable "torque_token" {
  type        = string
  sensitive   = true
  description = "Torque token (long token or short token) to be used for Torque API authneticaiton."
}

variable "environment_name" {
  type        = string
  description = "Name of the environment"
}

variable "blueprint_name" {
  type        = string
  description = "Name of the blueprint"
}

variable "duration" {
  type        = string
  description = "Duration of the environment"
}

variable "space" {
  type        = string
  description = "Space where the environment resides"
}

variable "inputs" {
  type        = map(string)
  description = "Input variables for the environment"
}

variable "owner_email" {
  type        = string
  description = "Email of the environment owner"
}

variable "automation" {
  type        = bool
  description = "Whether automation is enabled for the environment"
}

variable "description" {
  type        = string
  description = "Description of the environment"
}

variable "tags" {
  type        = map(string)
  description = "Tags associated with the environment"
  default     = {}
}

variable "workflows" {
  type = list(object({
    name = string
    schedules = list(object({
      scheduler  = string
      overridden = bool
    }))
    reminder         = number
    inputs_overrides = map(string)
  }))
  description = "Workflows associated with the environment"
  default     = []
}

variable "collaborators" {
  type = object({
    collaborators_emails = list(string)
    all_space_members    = bool
  })
  default = null
}

variable "blueprint_source" {
  type = object({
    blueprint_name  = string
    repository_name = string
    branch          = string
    commit          = string
  })
  default = null
}
