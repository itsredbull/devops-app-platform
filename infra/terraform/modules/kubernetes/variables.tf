variable "name_prefix" {
  type = string
}

variable "kubernetes_version" {
  type    = string
  default = "1.30"
}

variable "cluster_role_arn" {
  description = "IAM role ARN for EKS cluster"
  type        = string
}

variable "node_role_arn" {
  description = "IAM role ARN for EKS nodes"
  type        = string
}

variable "subnet_ids" {
  type = list(string)
}

variable "node_instance_types" {
  type    = list(string)
  default = ["t3.medium"]
}

variable "desired_nodes" {
  type    = number
  default = 2
}

variable "min_nodes" {
  type    = number
  default = 1
}

variable "max_nodes" {
  type    = number
  default = 3
}

variable "tags" {
  type    = map(string)
  default = {}
}
