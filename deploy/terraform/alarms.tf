# CloudWatch Alarms for RDS

resource "aws_cloudwatch_metric_alarm" "rds_cpu" {
  count = var.enable_alarms ? 1 : 0

  alarm_name          = "augment-fund-${var.environment}-rds-cpu-high"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 2
  metric_name         = "CPUUtilization"
  namespace           = "AWS/RDS"
  period              = 300
  statistic           = "Average"
  threshold           = 80
  alarm_description   = "RDS CPU utilization is above 80%"

  dimensions = {
    DBInstanceIdentifier = module.rds.db_instance_id
  }

  tags = {
    Name = "augment-fund-${var.environment}-rds-cpu-alarm"
  }
}

resource "aws_cloudwatch_metric_alarm" "rds_connections" {
  count = var.enable_alarms ? 1 : 0

  alarm_name          = "augment-fund-${var.environment}-rds-connections-high"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 2
  metric_name         = "DatabaseConnections"
  namespace           = "AWS/RDS"
  period              = 300
  statistic           = "Average"
  threshold           = 80
  alarm_description   = "RDS connections are above 80"

  dimensions = {
    DBInstanceIdentifier = module.rds.db_instance_id
  }

  tags = {
    Name = "augment-fund-${var.environment}-rds-connections-alarm"
  }
}

resource "aws_cloudwatch_metric_alarm" "rds_storage" {
  count = var.enable_alarms ? 1 : 0

  alarm_name          = "augment-fund-${var.environment}-rds-storage-low"
  comparison_operator = "LessThanThreshold"
  evaluation_periods  = 2
  metric_name         = "FreeStorageSpace"
  namespace           = "AWS/RDS"
  period              = 300
  statistic           = "Average"
  threshold           = 5368709120 # 5 GB in bytes
  alarm_description   = "RDS free storage is below 5 GB"

  dimensions = {
    DBInstanceIdentifier = module.rds.db_instance_id
  }

  tags = {
    Name = "augment-fund-${var.environment}-rds-storage-alarm"
  }
}

# CloudWatch Alarms for ECS

resource "aws_cloudwatch_metric_alarm" "ecs_cpu" {
  count = var.enable_alarms ? 1 : 0

  alarm_name          = "augment-fund-${var.environment}-ecs-cpu-high"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 2
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = 300
  statistic           = "Average"
  threshold           = 80
  alarm_description   = "ECS CPU utilization is above 80%"

  dimensions = {
    ClusterName = module.ecs.cluster_name
    ServiceName = module.ecs.service_name
  }

  tags = {
    Name = "augment-fund-${var.environment}-ecs-cpu-alarm"
  }
}

resource "aws_cloudwatch_metric_alarm" "ecs_memory" {
  count = var.enable_alarms ? 1 : 0

  alarm_name          = "augment-fund-${var.environment}-ecs-memory-high"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 2
  metric_name         = "MemoryUtilization"
  namespace           = "AWS/ECS"
  period              = 300
  statistic           = "Average"
  threshold           = 80
  alarm_description   = "ECS memory utilization is above 80%"

  dimensions = {
    ClusterName = module.ecs.cluster_name
    ServiceName = module.ecs.service_name
  }

  tags = {
    Name = "augment-fund-${var.environment}-ecs-memory-alarm"
  }
}
