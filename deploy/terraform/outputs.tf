# Secrets outputs
output "db_credentials_secret_arn" {
  description = "ARN of the database credentials secret"
  value       = aws_secretsmanager_secret.db_credentials.arn
}

output "db_credentials_secret_name" {
  description = "Name of the database credentials secret"
  value       = aws_secretsmanager_secret.db_credentials.name
}

output "secrets_access_policy_arn" {
  description = "ARN of the IAM policy for secrets access"
  value       = aws_iam_policy.secrets_access.arn
}

# VPC outputs
output "vpc_id" {
  description = "VPC ID"
  value       = module.vpc.vpc_id
}

# RDS outputs
output "rds_endpoint" {
  description = "RDS endpoint"
  value       = module.rds.endpoint
}

output "rds_address" {
  description = "RDS address"
  value       = module.rds.address
}

# ECR outputs
output "ecr_repository_url" {
  description = "ECR repository URL"
  value       = module.ecr.repository_url
}

# ECS outputs
output "ecs_cluster_name" {
  description = "ECS cluster name"
  value       = module.ecs.cluster_name
}

output "api_url" {
  description = "API URL (via ALB)"
  value       = module.ecs.api_url
}

output "alb_dns_name" {
  description = "ALB DNS name"
  value       = module.ecs.alb_dns_name
}

# Frontend outputs
output "frontend_bucket_name" {
  description = "Frontend S3 bucket name"
  value       = module.frontend.bucket_name
}

output "frontend_website_url" {
  description = "Frontend website URL"
  value       = module.frontend.website_url
}
