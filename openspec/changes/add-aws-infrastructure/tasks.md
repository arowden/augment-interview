## 1. Root Configuration
- [ ] 1.1 Create terraform/main.tf with provider and module composition
- [ ] 1.2 Create terraform/variables.tf with input variables
- [ ] 1.3 Create terraform/outputs.tf with important outputs
- [ ] 1.4 Create terraform/terraform.tfvars.example

## 2. VPC Module
- [ ] 2.1 Create terraform/modules/vpc/main.tf
- [ ] 2.2 Create VPC with DNS support
- [ ] 2.3 Create public and private subnets in 2 AZs
- [ ] 2.4 Create internet gateway (no NAT Gateway - using VPC endpoints)
- [ ] 2.5 Create route tables
- [ ] 2.6 Create security groups (ALB, ECS, RDS, VPC Endpoints)
- [ ] 2.7 Create VPC endpoint for ECR Docker (interface)
- [ ] 2.8 Create VPC endpoint for ECR API (interface)
- [ ] 2.9 Create VPC endpoint for S3 (gateway)
- [ ] 2.10 Create VPC endpoint for CloudWatch Logs (interface)
- [ ] 2.11 Enable VPC flow logs for REJECT traffic
- [ ] 2.12 Create flow logs IAM role and CloudWatch log group

## 3. Secrets Module
- [ ] 3.1 Create terraform/modules/secrets/main.tf
- [ ] 3.2 Create Secrets Manager secret for DB password
- [ ] 3.3 Generate random password
- [ ] 3.4 Set recovery_window_in_days to 7
- [ ] 3.5 Output secret ARN for ECS task

## 4. RDS Module with Security
- [ ] 4.1 Create terraform/modules/rds/main.tf
- [ ] 4.2 Create KMS key for RDS encryption
- [ ] 4.3 Create RDS PostgreSQL instance with storage_encrypted = true
- [ ] 4.4 Configure backup_retention_period = 7
- [ ] 4.5 Configure deletion_protection = true
- [ ] 4.6 Create parameter group with log_min_duration_statement = 1000
- [ ] 4.7 Create DB subnet group
- [ ] 4.8 Configure security group for private access
- [ ] 4.9 Reference password from Secrets Manager

## 5. ECR Module
- [ ] 5.1 Create terraform/modules/ecr/main.tf
- [ ] 5.2 Create ECR repository for API image
- [ ] 5.3 Configure lifecycle policy for image cleanup
- [ ] 5.4 Output repository URL and ARN

## 6. ECS Module with HTTPS
- [ ] 6.1 Create terraform/modules/ecs/main.tf
- [ ] 6.2 Create ECS cluster
- [ ] 6.3 Create task definition with secrets reference
- [ ] 6.4 Create ECS service with desired_count = 2
- [ ] 6.5 Create ALB with target group
- [ ] 6.6 Create ACM certificate for HTTPS
- [ ] 6.7 Create HTTPS listener on port 443 with TLS 1.3
- [ ] 6.8 Create HTTP listener on port 80 with redirect to HTTPS
- [ ] 6.9 Create scoped IAM execution role (ECR repo ARN, log group ARN, secret ARN)
- [ ] 6.10 Create task role
- [ ] 6.11 Configure CloudWatch log group with 7-day retention

## 7. CloudWatch Alarms
- [ ] 7.1 Create alarm for RDS CPU > 80%
- [ ] 7.2 Create alarm for RDS connections > 80
- [ ] 7.3 Configure alarm actions (SNS topic optional)

## 8. Frontend Module with Encryption
- [ ] 8.1 Create terraform/modules/frontend/main.tf
- [ ] 8.2 Create S3 bucket with unique name
- [ ] 8.3 Configure SSE-S3 encryption
- [ ] 8.4 Configure website hosting
- [ ] 8.5 Configure bucket policy for public read
- [ ] 8.6 Output website URL

## 9. Documentation
- [ ] 9.1 Create deployment instructions in README
- [ ] 9.2 Document required AWS permissions
- [ ] 9.3 Document ACM certificate setup
- [ ] 9.4 Document cost estimate (~$97/month)
