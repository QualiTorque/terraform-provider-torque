variable "torque_space" {
  type        = string
  description = "Torque space to be used during Torque API authentication."
}

variable "torque_token" {
  type        = string
  sensitive   = true
  description = "Torque token (long token or short token) to be used for Torque API authneticaiton."
}

variable "torque_host" {
  type        = string
  default     = "https://portal.qtorque.io/"
  description = "Torque portal URL"
}

variable "space_name" {
  type        = string
  description = "Torque space where the new repository should be onboarded in"
}

variable "repository_url" {
  type        = string
  description = "The full URL of the repository to be used"
}

variable "token" {
  type        = string
  sensitive   = true
  description = "The access token to the repository"
}

variable "branch" {
  type        = string
  description = "The repository branch to be used"
}

variable "repository_name" {
  type        = string
  description = "The logical name of the repository to be used in Torque"
}

variable "base_url" {
  type        = string
  description = "The base URL of the Gitlab Enetperise instance to be used in Torque"
}

variable "use_all_agents" {
  type        = bool
  description = "Whether all associated agents can be used to onboard and sync this repository. Must be set to false if agents attribute is used."
}

variable "agents" {
  type        = string
  description = "List of specific agents to use to onboard and sync this repository. Cannot be specified when use_all_agents is true."
}
