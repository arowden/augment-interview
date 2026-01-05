# Change: Add AWS Infrastructure

## Why
Production deployment requires cloud infrastructure. Terraform enables reproducible, version-controlled infrastructure as code. A simplified setup (static IP, no HTTPS) reduces complexity while meeting demo requirements.

## What Changes
- Add Terraform modules for VPC, RDS, ECS, ECR, and S3
- Configure ECS Fargate for running Go API containers
- Configure RDS PostgreSQL for database
- Configure S3 for static frontend hosting
- Output static IP and URLs for access

## Impact
- Affected specs: aws-infrastructure (new)
- Affected code: `/terraform/`
