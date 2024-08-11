variable "torque_space" {
  type        = string
  description = "Torque space to be used during Torque API authentication."
}

variable "torque_token" {
  type        = string
  sensitive   = true
  description = "Torque token (long token or short token) to be used for Torque API authneticaiton."
}

variable "resource_name" {
  type        = string
  description = "Resource name to be presented in the resouce catalog card"
}

variable "resource_image" {
  type        = string
  default     = "https://portal.qtorque.io/static/media/networking.dc1b7fb73182de0136d059a1eb00af4f.svg"
  description = "Resource image to be presented in the resouce catalog card. Image should be uploaded to Torque prior to usage."
}

variable "resource_data" {
  description = "Resource data to be oresented in the resouce catalog card"
  type        = map(string)
}
