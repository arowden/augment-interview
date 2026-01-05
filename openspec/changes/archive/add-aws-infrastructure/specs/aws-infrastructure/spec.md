## ADDED Requirements

### Requirement: Terraform Project Structure
The infrastructure SHALL be defined as Terraform modules in the /terraform directory.

#### Scenario: Root configuration
- **WHEN** terraform/ is examined
- **THEN** it contains main.tf, variables.tf, outputs.tf, and modules/ directory

#### Scenario: Module structure
- **WHEN** terraform/modules/ is examined
- **THEN** it contains vpc, rds, ecr, ecs, and frontend subdirectories

#### Scenario: Variable file template
- **WHEN** terraform.tfvars.example is examined
- **THEN** it shows all required variables with placeholder values

### Requirement: AWS Provider Configuration
The system SHALL configure the AWS provider with region and credentials.

#### Scenario: Provider configuration
- **WHEN** terraform init is run
- **THEN** AWS provider is downloaded and configured

#### Scenario: Region configuration
- **WHEN** aws_region variable is set to "us-east-1"
- **THEN** all resources are created in us-east-1

### Requirement: VPC Module
The system SHALL create a VPC with public and private subnets.

#### Scenario: VPC creation
- **WHEN** vpc module is applied
- **THEN** a VPC with CIDR 10.0.0.0/16 is created

#### Scenario: Subnet creation
- **WHEN** vpc module is applied
- **THEN** 2 public subnets and 2 private subnets are created across 2 AZs

#### Scenario: Internet gateway
- **WHEN** vpc module is applied
- **THEN** an internet gateway is attached to the VPC

#### Scenario: NAT gateway
- **WHEN** vpc module is applied
- **THEN** a NAT gateway is created in a public subnet

#### Scenario: Route tables
- **WHEN** vpc module is applied
- **THEN** public subnets route to IGW, private subnets route to NAT

### Requirement: Security Groups
The VPC module SHALL create security groups for ALB, ECS, RDS, and VPC Endpoints.

#### Scenario: ALB security group
- **WHEN** alb security group is examined
- **THEN** it allows inbound TCP 443 and 80 from 0.0.0.0/0

#### Scenario: ECS security group
- **WHEN** ecs security group is examined
- **THEN** it allows inbound TCP 8080 from ALB security group only

#### Scenario: RDS security group
- **WHEN** rds security group is examined
- **THEN** it allows inbound TCP 5432 from ECS security group only

#### Scenario: VPC Endpoints security group
- **WHEN** vpce security group is examined
- **THEN** it allows inbound TCP 443 from ECS security group

### Requirement: RDS Module with Encryption and Backups
The system SHALL create an encrypted PostgreSQL RDS instance with automated backups.

#### Scenario: RDS instance creation
- **WHEN** rds module is applied
- **THEN** a db.t3.micro PostgreSQL 16 instance is created

#### Scenario: Storage encryption
- **WHEN** RDS instance is examined
- **THEN** storage_encrypted is true with KMS key

#### Scenario: Automated backups
- **WHEN** RDS instance is examined
- **THEN** backup_retention_period is 7 days

#### Scenario: Deletion protection
- **WHEN** RDS instance is examined
- **THEN** deletion_protection is true

#### Scenario: Parameter group for logging
- **WHEN** RDS parameter group is examined
- **THEN** log_min_duration_statement is set to 1000ms

#### Scenario: DB subnet group
- **WHEN** rds module is applied
- **THEN** a DB subnet group using private subnets is created

#### Scenario: Private access only
- **WHEN** RDS instance is examined
- **THEN** publicly_accessible is false

### Requirement: ECR Module
The system SHALL create an ECR repository for the API container image.

#### Scenario: Repository creation
- **WHEN** ecr module is applied
- **THEN** an ECR repository named "augment-fund-api" is created

#### Scenario: Lifecycle policy
- **WHEN** ecr repository is examined
- **THEN** a lifecycle policy keeps only last 10 images

#### Scenario: Repository URL output
- **WHEN** ecr module outputs are examined
- **THEN** repository_url output is available for docker push

### Requirement: ECS Module
The system SHALL create an ECS Fargate service running the API container.

#### Scenario: Cluster creation
- **WHEN** ecs module is applied
- **THEN** an ECS cluster is created

#### Scenario: Task definition
- **WHEN** ecs module is applied
- **THEN** a task definition with 0.5 vCPU and 1GB memory is created

#### Scenario: Container configuration
- **WHEN** task definition is examined
- **THEN** it uses the ECR image and exposes port 8080

#### Scenario: Service creation
- **WHEN** ecs module is applied
- **THEN** an ECS service with desired count 1 is created in private subnets

#### Scenario: Environment variables
- **WHEN** task definition is examined
- **THEN** DATABASE_URL environment variable is configured

### Requirement: Application Load Balancer with HTTPS
The ECS module SHALL create an ALB with HTTPS and HTTP redirect.

#### Scenario: ALB creation
- **WHEN** ecs module is applied
- **THEN** an internet-facing ALB is created in public subnets

#### Scenario: Target group
- **WHEN** ecs module is applied
- **THEN** a target group with health check on /health is created

#### Scenario: HTTPS Listener
- **WHEN** ecs module is applied
- **THEN** an HTTPS listener on port 443 forwards to target group with TLS 1.3

#### Scenario: HTTP to HTTPS Redirect
- **WHEN** ecs module is applied
- **THEN** an HTTP listener on port 80 redirects to HTTPS (301)

#### Scenario: ACM Certificate
- **WHEN** ecs module is applied
- **THEN** an ACM certificate is provisioned for HTTPS

### Requirement: IAM Roles with Scoped Permissions
The ECS module SHALL create IAM roles with least-privilege permissions.

#### Scenario: Execution role with scoped ECR
- **WHEN** execution role policy is examined
- **THEN** ECR permissions are scoped to specific repository ARN

#### Scenario: Execution role with scoped logs
- **WHEN** execution role policy is examined
- **THEN** CloudWatch Logs permissions are scoped to specific log group ARN

#### Scenario: Execution role with secrets access
- **WHEN** execution role policy is examined
- **THEN** Secrets Manager permissions are scoped to specific secret ARN

#### Scenario: Task role
- **WHEN** ecs module is applied
- **THEN** a task role is created for runtime permissions

### Requirement: CloudWatch Logging
The ECS module SHALL configure CloudWatch logs for container output.

#### Scenario: Log group creation
- **WHEN** ecs module is applied
- **THEN** a log group /ecs/augment-fund-api is created

#### Scenario: Log retention
- **WHEN** log group is examined
- **THEN** retention is set to 7 days

#### Scenario: Container logging
- **WHEN** container runs
- **THEN** stdout/stderr are sent to CloudWatch logs

### Requirement: Frontend Module with Encryption
The system SHALL create an encrypted S3 bucket for static website hosting.

#### Scenario: Bucket creation
- **WHEN** frontend module is applied
- **THEN** an S3 bucket is created with unique name

#### Scenario: S3 encryption
- **WHEN** bucket is examined
- **THEN** server-side encryption (SSE-S3) is enabled

#### Scenario: Website configuration
- **WHEN** bucket is examined
- **THEN** static website hosting is enabled with index.html as index document

#### Scenario: Public access
- **WHEN** bucket policy is examined
- **THEN** it allows public read access to all objects

#### Scenario: Website URL output
- **WHEN** frontend module outputs are examined
- **THEN** website_url output is available

### Requirement: Secrets Manager Module
The system SHALL store sensitive values in AWS Secrets Manager.

#### Scenario: Database password secret
- **WHEN** secrets module is applied
- **THEN** a secret named "augment-fund/db-password" is created

#### Scenario: Secret value generation
- **WHEN** secret is created
- **THEN** a random password is generated and stored

#### Scenario: Recovery window
- **WHEN** secret is examined
- **THEN** recovery_window_in_days is set to 7

### Requirement: VPC Endpoints
The system SHALL create VPC endpoints to avoid NAT Gateway costs.

#### Scenario: ECR Docker endpoint
- **WHEN** vpc module is applied
- **THEN** an interface endpoint for ecr.dkr is created

#### Scenario: ECR API endpoint
- **WHEN** vpc module is applied
- **THEN** an interface endpoint for ecr.api is created

#### Scenario: S3 endpoint
- **WHEN** vpc module is applied
- **THEN** a gateway endpoint for S3 is created

### Requirement: VPC Flow Logs
The system SHALL enable VPC flow logs for security monitoring.

#### Scenario: Flow logs enabled
- **WHEN** vpc module is applied
- **THEN** flow logs are configured for REJECT traffic

#### Scenario: Flow logs destination
- **WHEN** flow logs are examined
- **THEN** they are sent to CloudWatch Logs

### Requirement: CloudWatch Alarms
The system SHALL create CloudWatch alarms for critical metrics.

#### Scenario: RDS CPU alarm
- **WHEN** ecs module is applied
- **THEN** an alarm for RDS CPU > 80% is created

#### Scenario: RDS connections alarm
- **WHEN** ecs module is applied
- **THEN** an alarm for RDS connections > 80 is created

### Requirement: ECS High Availability
The ECS service SHALL run multiple tasks for availability.

#### Scenario: ECS desired count
- **WHEN** ecs module is applied
- **THEN** desired_count is set to 2 minimum

### Requirement: Terraform Outputs
The root module SHALL output all necessary URLs and values for deployment.

#### Scenario: API URL output
- **WHEN** terraform output is run
- **THEN** api_url shows the ALB URL with Elastic IP

#### Scenario: Frontend URL output
- **WHEN** terraform output is run
- **THEN** frontend_url shows the S3 website URL

#### Scenario: ECR URL output
- **WHEN** terraform output is run
- **THEN** ecr_url shows the repository URL for docker push

#### Scenario: Database URL output
- **WHEN** terraform output is run
- **THEN** database_url shows the RDS connection string (sensitive)

### Requirement: Input Variables
The root module SHALL accept configurable input variables.

#### Scenario: Required variables
- **WHEN** terraform plan is run without variables
- **THEN** it prompts for aws_region, db_password

#### Scenario: Default values
- **WHEN** variable definitions are examined
- **THEN** aws_region defaults to "us-east-1"

#### Scenario: Sensitive variables
- **WHEN** db_password variable is defined
- **THEN** it is marked as sensitive
