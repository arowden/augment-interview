## 1. Root Configuration
- [x] 1.1 Create terraform/main.tf with provider and module composition
- [x] 1.2 Create terraform/variables.tf with input variables
- [x] 1.3 Create terraform/outputs.tf with important outputs
- [x] 1.4 Create terraform/terraform.tfvars.example

## 2. VPC Module
- [x] 2.1 Create terraform/modules/vpc/main.tf
- [x] 2.2 Create VPC with DNS support
- [x] 2.3 Create public and private subnets in 2 AZs
- [x] 2.4 Create internet gateway (no NAT Gateway - using VPC endpoints)
- [x] 2.5 Create route tables
- [x] 2.6 Create security groups (ALB, ECS, RDS, VPC Endpoints)
- [x] 2.7 Create VPC endpoint for ECR Docker (interface)
- [x] 2.8 Create VPC endpoint for ECR API (interface)
- [x] 2.9 Create VPC endpoint for S3 (gateway)
- [x] 2.10 Create VPC endpoint for CloudWatch Logs (interface)
- [x] 2.11 Enable VPC flow logs for REJECT traffic
- [x] 2.12 Create flow logs IAM role and CloudWatch log group

## 3. Secrets Module
- [x] 3.1 Create terraform/secrets.tf (in root, not module)
- [x] 3.2 Create Secrets Manager secret for DB password
- [x] 3.3 Generate random password
- [x] 3.4 Set recovery_window_in_days to 7
- [x] 3.5 Output secret ARN for ECS task

## 4. RDS Module with Security
- [x] 4.1 Create terraform/modules/rds/main.tf
- [x] 4.2 Create KMS key for RDS encryption
- [x] 4.3 Create RDS PostgreSQL instance with storage_encrypted = true
- [x] 4.4 Configure backup_retention_period = 7
- [x] 4.5 Configure deletion_protection = true
- [x] 4.6 Create parameter group with log_min_duration_statement = 1000
- [x] 4.7 Create DB subnet group
- [x] 4.8 Configure security group for private access
- [x] 4.9 Reference password from Secrets Manager

## 5. ECR Module
- [x] 5.1 Create terraform/modules/ecr/main.tf
- [x] 5.2 Create ECR repository for API image
- [x] 5.3 Configure lifecycle policy for image cleanup
- [x] 5.4 Output repository URL and ARN

## 6. ECS Module with HTTPS
- [x] 6.1 Create terraform/modules/ecs/main.tf
- [x] 6.2 Create ECS cluster
- [x] 6.3 Create task definition with secrets reference
- [x] 6.4 Create ECS service with desired_count = 2
- [x] 6.5 Create ALB with target group
- [x] 6.6 Create ACM certificate for HTTPS (optional, domain-based)
- [x] 6.7 Create HTTPS listener on port 443 with TLS 1.3 (when domain provided)
- [x] 6.8 Create HTTP listener on port 80 with redirect to HTTPS (or direct access)
- [x] 6.9 Create scoped IAM execution role (ECR repo ARN, log group ARN, secret ARN)
- [x] 6.10 Create task role
- [x] 6.11 Configure CloudWatch log group with 7-day retention

## 7. CloudWatch Alarms
- [x] 7.1 Create alarm for RDS CPU > 80%
- [x] 7.2 Create alarm for RDS connections > 80
- [x] 7.3 Create alarm for RDS storage < 5GB
- [x] 7.4 Create alarm for ECS CPU > 80%
- [x] 7.5 Create alarm for ECS memory > 80%
- [x] 7.6 Configure alarm toggle via enable_alarms variable

## 8. Frontend Module with Encryption
- [x] 8.1 Create terraform/modules/frontend/main.tf
- [x] 8.2 Create S3 bucket with unique name
- [x] 8.3 Configure SSE-S3 encryption
- [x] 8.4 Configure website hosting
- [x] 8.5 Configure bucket policy for public read
- [x] 8.6 Output website URL

## 9. Documentation
- [x] 9.1 Create deployment instructions in README
- [x] 9.2 Document required AWS permissions
- [x] 9.3 Document ACM certificate setup
- [x] 9.4 Document cost estimate (~$97/month)
