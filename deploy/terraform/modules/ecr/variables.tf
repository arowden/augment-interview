variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "image_retention_count" {
  description = "Number of images to retain"
  type        = number
  default     = 10
}
