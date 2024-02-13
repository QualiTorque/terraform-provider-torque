variable "torque_space" {
  type        = string
  description = "Torque space to be used during Torque API authentication."
}

variable "torque_token" {
  type        = string
  sensitive   = true
  description = "Torque token (long token or short token) to be used for Torque API authneticaiton."
}

variable "space_name" {
  type        = string
  description = "Torque space where the new repository should be onboarded in"
}

variable "repository_url" {
  type        = string
  description = "The URL of the repository to be used"
}

variable "access_token" {
  type        = string
  sensitive   = true
  description = "The access token to the repository"
}

variable "repository_type" {
  type        = string
  description = "The repository type to use. Possible values:  Available types: github, bitbucket, gitlab, azure (for Azure DevOps)"
}

variable "branch" {
  type        = string
  description = "The repository branch to be used"
}

variable "repository_name" {
  type        = string
  description = "The logical name of the repository to be used in Torque"
}