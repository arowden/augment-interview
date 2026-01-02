## Context
AWS infrastructure hosts the production deployment. The architecture balances security best practices with cost-effectiveness. Terraform modules provide reusable, composable infrastructure components.

## Goals / Non-Goals
- Goals: Secure production environment, reproducible infrastructure, cost-effective, encryption at rest
- Non-Goals: High availability (multi-region), auto-scaling, CI/CD pipeline

## Decisions
- Decision: HTTPS via ACM certificate (self-managed for demo, or domain-validated)
- Alternatives considered: HTTP only (security risk for credentials), custom certs (management overhead)

- Decision: VPC Endpoints for ECR (saves NAT Gateway costs)
- Alternatives considered: NAT Gateway (~$35/month), public ECR pull (security risk)

- Decision: Secrets Manager for database password
- Alternatives considered: Environment variables (visible in console), SSM Parameter Store (less features)

- Decision: S3 static website hosting with SSE encryption
- Alternatives considered: Unencrypted (fails audit), CloudFront CDN (overkill)

- Decision: ECS Fargate with 2 tasks minimum (availability)
- Alternatives considered: Single task (no availability), EC2 (management overhead)

- Decision: RDS single-AZ with encryption and 7-day backups
- Alternatives considered: Multi-AZ (2x cost), no backups (data loss risk)

## Module Structure
```
terraform/
  main.tf                 # Provider, backend, module composition
  variables.tf            # Input variables
  outputs.tf              # Important outputs
  terraform.tfvars.example
  modules/
    vpc/
      main.tf             # VPC, subnets, gateways, security groups, flow logs
      variables.tf
      outputs.tf
    rds/
      main.tf             # PostgreSQL instance with encryption and backups
      variables.tf
      outputs.tf
    ecr/
      main.tf             # Container registry
      variables.tf
      outputs.tf
    ecs/
      main.tf             # Cluster, service, ALB with HTTPS, scoped IAM
      variables.tf
      outputs.tf
    frontend/
      main.tf             # S3 bucket with encryption
      variables.tf
      outputs.tf
    secrets/
      main.tf             # Secrets Manager for DB password
      variables.tf
      outputs.tf
```

## Network Design
```
VPC (10.0.0.0/16)
├── Public Subnet A (10.0.1.0/24) - ALB
├── Public Subnet B (10.0.2.0/24) - ALB
├── Private Subnet A (10.0.10.0/24) - ECS Tasks, RDS
├── Private Subnet B (10.0.11.0/24) - ECS Tasks, RDS
└── VPC Endpoints (ECR, S3, CloudWatch Logs)
```

## Security Groups
- **alb-sg**: Inbound 443 from 0.0.0.0/0, Inbound 80 from 0.0.0.0/0 (redirect)
- **ecs-sg**: Inbound 8080 from alb-sg only
- **rds-sg**: Inbound 5432 from ecs-sg only
- **vpce-sg**: Inbound 443 from ecs-sg (for VPC endpoints)

## RDS Configuration with Security
```hcl
resource "aws_db_instance" "main" {
  identifier             = "augment-fund-db"
  engine                 = "postgres"
  engine_version         = "16"
  instance_class         = "db.t3.micro"
  allocated_storage      = 20

  storage_encrypted      = true
  kms_key_id            = aws_kms_key.rds.arn

  backup_retention_period = 7
  backup_window          = "03:00-04:00"
  maintenance_window     = "Mon:04:00-Mon:05:00"

  deletion_protection    = true
  copy_tags_to_snapshot  = true
  skip_final_snapshot    = false
  final_snapshot_identifier = "augment-fund-final"

  publicly_accessible    = false
  db_subnet_group_name   = aws_db_subnet_group.main.name
  vpc_security_group_ids = [aws_security_group.rds.id]

  parameter_group_name   = aws_db_parameter_group.main.name
}

resource "aws_db_parameter_group" "main" {
  family = "postgres16"
  name   = "augment-fund-params"

  parameter {
    name  = "log_min_duration_statement"
    value = "1000"
  }
}
```

## S3 Configuration with Encryption
```hcl
resource "aws_s3_bucket_server_side_encryption_configuration" "frontend" {
  bucket = aws_s3_bucket.frontend.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}
```

## Secrets Manager for Database Password
```hcl
resource "aws_secretsmanager_secret" "db_password" {
  name = "augment-fund/db-password"

  recovery_window_in_days = 7
}

resource "aws_secretsmanager_secret_version" "db_password" {
  secret_id     = aws_secretsmanager_secret.db_password.id
  secret_string = random_password.db_password.result
}
```

## Scoped IAM Policy
```hcl
resource "aws_iam_role_policy" "ecs_execution" {
  name = "ecs-execution-policy"
  role = aws_iam_role.ecs_execution.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ecr:GetAuthorizationToken"
        ]
        Resource = "*"
      },
      {
        Effect = "Allow"
        Action = [
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage"
        ]
        Resource = [aws_ecr_repository.api.arn]
      },
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Resource = ["${aws_cloudwatch_log_group.ecs.arn}:*"]
      },
      {
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue"
        ]
        Resource = [aws_secretsmanager_secret.db_password.arn]
      }
    ]
  })
}
```

## HTTPS Configuration
```hcl
resource "aws_lb_listener" "https" {
  load_balancer_arn = aws_lb.main.arn
  port              = 443
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-TLS13-1-2-2021-06"
  certificate_arn   = aws_acm_certificate.main.arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.main.arn
  }
}

resource "aws_lb_listener" "http_redirect" {
  load_balancer_arn = aws_lb.main.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}
```

## VPC Endpoints (Cost Savings)
```hcl
resource "aws_vpc_endpoint" "ecr_dkr" {
  vpc_id              = aws_vpc.main.id
  service_name        = "com.amazonaws.${var.region}.ecr.dkr"
  vpc_endpoint_type   = "Interface"
  subnet_ids          = aws_subnet.private[*].id
  security_group_ids  = [aws_security_group.vpce.id]
  private_dns_enabled = true
}

resource "aws_vpc_endpoint" "ecr_api" {
  vpc_id              = aws_vpc.main.id
  service_name        = "com.amazonaws.${var.region}.ecr.api"
  vpc_endpoint_type   = "Interface"
  subnet_ids          = aws_subnet.private[*].id
  security_group_ids  = [aws_security_group.vpce.id]
  private_dns_enabled = true
}

resource "aws_vpc_endpoint" "s3" {
  vpc_id            = aws_vpc.main.id
  service_name      = "com.amazonaws.${var.region}.s3"
  vpc_endpoint_type = "Gateway"
  route_table_ids   = aws_route_table.private[*].id
}
```

## VPC Flow Logs
```hcl
resource "aws_flow_log" "main" {
  vpc_id                   = aws_vpc.main.id
  traffic_type             = "REJECT"
  log_destination_type     = "cloud-watch-logs"
  log_destination          = aws_cloudwatch_log_group.flow_logs.arn
  iam_role_arn            = aws_iam_role.flow_logs.arn
}
```

## CloudWatch Alarms
```hcl
resource "aws_cloudwatch_metric_alarm" "rds_cpu" {
  alarm_name          = "rds-cpu-high"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 2
  metric_name         = "CPUUtilization"
  namespace           = "AWS/RDS"
  period              = 300
  statistic           = "Average"
  threshold           = 80

  dimensions = {
    DBInstanceIdentifier = aws_db_instance.main.id
  }
}

resource "aws_cloudwatch_metric_alarm" "rds_connections" {
  alarm_name          = "rds-connections-high"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 2
  metric_name         = "DatabaseConnections"
  namespace           = "AWS/RDS"
  period              = 300
  statistic           = "Average"
  threshold           = 80

  dimensions = {
    DBInstanceIdentifier = aws_db_instance.main.id
  }
}
```

## ECS Task Definition with Secrets
```json
{
  "containerDefinitions": [{
    "name": "api",
    "image": "${ecr_url}:${version}",
    "portMappings": [{"containerPort": 8080}],
    "secrets": [
      {
        "name": "DATABASE_PASSWORD",
        "valueFrom": "${db_password_secret_arn}"
      }
    ],
    "environment": [
      {"name": "DATABASE_HOST", "value": "${rds_host}"},
      {"name": "DATABASE_NAME", "value": "augment_fund"},
      {"name": "DATABASE_USER", "value": "augment"},
      {"name": "OTEL_ENABLED", "value": "true"},
      {"name": "OTEL_EXPORTER_OTLP_ENDPOINT", "value": "http://localhost:4317"}
    ],
    "logConfiguration": {
      "logDriver": "awslogs",
      "options": {"awslogs-group": "/ecs/augment-fund-api"}
    }
  }]
}
```

## Cost Estimate (us-east-1) - Updated
- ECS Fargate (0.5 vCPU, 1GB) x2: ~$30/month
- RDS db.t3.micro with backups: ~$18/month
- ALB with HTTPS: ~$25/month
- S3: <$1/month
- VPC Endpoints (3): ~$22/month
- Secrets Manager: ~$0.50/month
- KMS: ~$1/month
- **Total**: ~$97/month

*Note: VPC endpoints ($22) vs NAT Gateway ($35) saves ~$13/month*

## Risks / Trade-offs
- Single-AZ RDS has downtime risk → Acceptable for demo, backups mitigate data loss
- VPC endpoints have hourly cost → Cheaper than NAT Gateway for this use case
- ACM certificate requires domain validation → Self-signed for development
- Two ECS tasks doubles compute cost → Required for zero-downtime deploys

## Open Questions
- None
