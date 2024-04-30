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
  description = "The space name where the notification should be created on."
}

variable "notification_name" {
  type        = string
  description = "The name of the notification to be created."
}

variable "environment_launched" {
  type        = bool
  default     = false
  description = "Get notifications when environments are launched"
}

variable "environment_deployed" {
  type        = bool
  default     = false
  description = "Get notifications when environments are deployed"
}

variable "environment_force_ended" {
  type        = bool
  default     = false
  description = "Get notifications when environments are forced ended"
}

variable "environment_idle" {
  type        = bool
  default     = false
  description = "Get notifications when environments are idle"
}

variable "environment_extended" {
  type        = bool
  default     = false
  description = "Get notifications when environments are extended"
}

variable "drift_detected" {
  type        = bool
  default     = false
  description = "Get notifications when a drift is detected"
}

variable "workflow_failed" {
  type        = bool
  default     = false
  description = "Get notifications when workflows failed to execute"
}

variable "workflow_started" {
  type        = bool
  default     = false
  description = "Get notifications when workflows started to execute"
}

variable "updates_detected" {
  type        = bool
  default     = false
  description = "Get notifications when environment updated are detected"
}

variable "collaborator_added" {
  type        = bool
  default     = false
  description = "Get notifications when environment collaborators are added"
}

variable "action_failed" {
  type        = bool
  default     = false
  description = "Get notifications when environment action failed to run"
}

variable "environment_ending_failed" {
  type        = bool
  default     = false
  description = "Get notifications when environment termination failed"
}

variable "environment_ended" {
  type        = bool
  default     = false
  description = "Get notifications when environment termination ended successfuly"
}

variable "environment_active_with_error" {
  type        = bool
  default     = false
  description = "Get notifications when environment setup ended with an error"
}

variable "idle_reminders" {
  type        = list(number)
  default     = false
  description = "Get notifications when environment setup ended with an error"
}