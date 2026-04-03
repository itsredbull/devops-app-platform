output "api_alarm_name" {
  value = aws_cloudwatch_metric_alarm.api_cpu_high.alarm_name
}

output "db_alarm_name" {
  value = aws_cloudwatch_metric_alarm.db_cpu_high.alarm_name
}
