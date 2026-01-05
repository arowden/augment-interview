# KMS Key for RDS encryption
resource "aws_kms_key" "rds" {
  description             = "KMS key for RDS encryption - ${var.environment}"
  deletion_window_in_days = 7
  enable_key_rotation     = true

  tags = {
    Name = "augment-fund-${var.environment}-rds-key"
  }
}

resource "aws_kms_alias" "rds" {
  name          = "alias/augment-fund-${var.environment}-rds"
  target_key_id = aws_kms_key.rds.key_id
}

# DB Subnet Group
resource "aws_db_subnet_group" "main" {
  name       = "augment-fund-${var.environment}"
  subnet_ids = var.subnet_ids

  tags = {
    Name = "augment-fund-${var.environment}-db-subnet-group"
  }
}

# DB Parameter Group
resource "aws_db_parameter_group" "main" {
  family = "postgres16"
  name   = "augment-fund-${var.environment}-params"

  parameter {
    name  = "log_min_duration_statement"
    value = "1000"
  }

  parameter {
    name  = "log_statement"
    value = "ddl"
  }

  tags = {
    Name = "augment-fund-${var.environment}-params"
  }
}

# Fetch password from Secrets Manager
data "aws_secretsmanager_secret_version" "db_password" {
  secret_id = var.db_password_secret_arn
}

locals {
  db_credentials = jsondecode(data.aws_secretsmanager_secret_version.db_password.secret_string)
}

# RDS PostgreSQL Instance
resource "aws_db_instance" "main" {
  identifier = "augment-fund-${var.environment}"

  engine         = "postgres"
  engine_version = "16"
  instance_class = var.instance_class

  allocated_storage     = var.allocated_storage
  max_allocated_storage = var.allocated_storage * 2
  storage_type          = "gp3"
  storage_encrypted     = true
  kms_key_id            = aws_kms_key.rds.arn

  db_name  = var.db_name
  username = var.db_username
  password = local.db_credentials.password

  db_subnet_group_name   = aws_db_subnet_group.main.name
  vpc_security_group_ids = [var.security_group_id]
  parameter_group_name   = aws_db_parameter_group.main.name

  publicly_accessible = false
  multi_az            = false

  backup_retention_period = var.backup_retention_period
  backup_window           = "03:00-04:00"
  maintenance_window      = "Mon:04:00-Mon:05:00"

  deletion_protection       = var.environment == "prod" ? true : false
  skip_final_snapshot       = var.environment == "prod" ? false : true
  final_snapshot_identifier = var.environment == "prod" ? "augment-fund-${var.environment}-final" : null
  copy_tags_to_snapshot     = true

  performance_insights_enabled = false

  tags = {
    Name = "augment-fund-${var.environment}-db"
  }
}
