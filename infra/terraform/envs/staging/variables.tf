variable "aws_region" {
  type = string
}

variable "name_prefix" {
  description = "Prefix for all resources"
  type        = string
  default     = "uptime-staging"
}

variable "vpc_cidr" {
  type    = string
  default = "10.20.0.0/16"
}

variable "public_subnet_cidrs" {
  type    = list(string)
  default = ["10.20.1.0/24", "10.20.2.0/24"]
}

variable "private_subnet_cidrs" {
  type    = list(string)
  default = ["10.20.11.0/24", "10.20.12.0/24"]
}

variable "availability_zones" {
  type    = list(string)
  default = ["us-east-1a", "us-east-1b"]
}

variable "enable_nat_gateway" {
  description = "Enable NAT for private subnet egress"
  type        = bool
  default     = true
}

variable "eks_cluster_role_arn" {
  description = "Existing IAM role ARN for EKS control plane"
  type        = string
}

variable "eks_node_role_arn" {
  description = "Existing IAM role ARN for EKS worker nodes"
  type        = string
}

variable "db_name" {
  type    = string
  default = "uptime"
}

variable "db_username" {
  type    = string
  default = "uptime"
}

variable "db_password" {
  type      = string
  sensitive = true
}

variable "alarm_actions" {
  description = "SNS topic ARNs for alarms"
  type        = list(string)
  default     = []
}

variable "tags" {
  type = map(string)
  default = {
    Project     = "devops-app-platform"
    Environment = "staging"
    ManagedBy   = "terraform"
  }
}
