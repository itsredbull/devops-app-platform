output "vpc_id" {
  value = module.network.vpc_id
}

output "eks_cluster_name" {
  value = module.kubernetes.cluster_name
}

output "db_endpoint" {
  value = module.database.db_endpoint
}

output "db_port" {
  value = module.database.db_port
}

output "monitoring_alarms" {
  value = {
    api = module.monitoring.api_alarm_name
    db  = module.monitoring.db_alarm_name
  }
}
