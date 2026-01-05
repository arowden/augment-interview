variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "bucket_name_suffix" {
  description = "Suffix for bucket name (to ensure uniqueness)"
  type        = string
  default     = ""
}
