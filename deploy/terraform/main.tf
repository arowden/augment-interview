# Main Terraform configuration for Augment Fund infrastructure

# VPC Module
module "vpc" {
  source = "./modules/vpc"

  environment        = var.environment
  aws_region         = var.aws_region
  vpc_cidr           = var.vpc_cidr
  availability_zones = var.availability_zones
}

# ECR Module
module "ecr" {
  source = "./modules/ecr"

  environment = var.environment
}

# RDS Module
module "rds" {
  source = "./modules/rds"

  environment            = var.environment
  vpc_id                 = module.vpc.vpc_id
  subnet_ids             = module.vpc.private_subnet_ids
  security_group_id      = module.vpc.rds_security_group_id
  db_name                = var.db_name
  db_username            = var.db_username
  db_password_secret_arn = aws_secretsmanager_secret.db_credentials.arn

  depends_on = [aws_secretsmanager_secret_version.db_credentials]
}

# ECS Module
module "ecs" {
  source = "./modules/ecs"

  environment            = var.environment
  aws_region             = var.aws_region
  vpc_id                 = module.vpc.vpc_id
  public_subnet_ids      = module.vpc.public_subnet_ids
  private_subnet_ids     = module.vpc.private_subnet_ids
  alb_security_group_id  = module.vpc.alb_security_group_id
  ecs_security_group_id  = module.vpc.ecs_security_group_id
  ecr_repository_url     = module.ecr.repository_url
  ecr_repository_arn     = module.ecr.repository_arn
  db_host                = module.rds.address
  db_name                = var.db_name
  db_username            = var.db_username
  db_password_secret_arn = aws_secretsmanager_secret.db_credentials.arn
  domain_name            = var.domain_name
  cpu                    = var.ecs_cpu
  memory                 = var.ecs_memory
  desired_count          = var.ecs_desired_count

  depends_on = [module.rds]
}

# Frontend Module
module "frontend" {
  source = "./modules/frontend"

  environment = var.environment
}
