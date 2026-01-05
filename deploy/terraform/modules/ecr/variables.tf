variable "environment" {
  description = "Environment name"
  type        = string
}

variable "image_retention_count" {
  description = "Number of images to retain"
  type        = number
  default     = 10
}
