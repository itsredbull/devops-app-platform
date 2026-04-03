resource "aws_cloudwatch_metric_alarm" "api_cpu_high" {
  alarm_name          = "${var.name_prefix}-api-cpu-high"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 2
  metric_name         = "CPUUtilization"
  namespace           = "AWS/EKS"
  period              = 300
  statistic           = "Average"
  threshold           = var.api_cpu_threshold
  alarm_description   = "EKS/API CPU high for uptime platform"
  treat_missing_data  = "notBreaching"

  dimensions = {
    ClusterName = var.cluster_name
  }

  alarm_actions = var.alarm_actions
  ok_actions    = var.alarm_actions

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-api-cpu-high"
  })
}

resource "aws_cloudwatch_metric_alarm" "db_cpu_high" {
  alarm_name          = "${var.name_prefix}-db-cpu-high"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 2
  metric_name         = "CPUUtilization"
  namespace           = "AWS/RDS"
  period              = 300
  statistic           = "Average"
  threshold           = var.db_cpu_threshold
  alarm_description   = "RDS CPU high for uptime platform"
  treat_missing_data  = "notBreaching"

  dimensions = {
    DBInstanceIdentifier = var.db_identifier
  }

  alarm_actions = var.alarm_actions
  ok_actions    = var.alarm_actions

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-db-cpu-high"
  })
}
