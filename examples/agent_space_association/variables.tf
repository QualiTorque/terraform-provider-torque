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
  description = "The space to be associated with the agent"
}

variable "agent_name" {
  type        = string
  description = "The agent name to be associated with the space"
}

variable "service_account" {
  type        = string
  description = "Kubernetes service account to be used by the agent for the specific space"
}

variable "namespace" {
  type        = string
  description = "Kubernetes namespace to be used by the agent for the specific space"
}
