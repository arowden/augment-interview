variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "dev"
}

variable "db_username" {
  description = "Database username"
  type        = string
  default     = "augment_fund"
}

variable "db_name" {
  description = "Database name"
  type        = string
  default     = "augment_fund"
}
