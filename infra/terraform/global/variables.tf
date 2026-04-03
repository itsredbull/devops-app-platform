variable "aws_region" {
  description = "AWS region for the remote state resources"
  type        = string
}

variable "state_bucket_name" {
  description = "Globally unique S3 bucket name for Terraform state"
  type        = string
}

variable "lock_table_name" {
  description = "DynamoDB table name used for state locking"
  type        = string
  default     = "devops-app-platform-tf-locks"
}

variable "tags" {
  description = "Tags applied to state resources"
  type        = map(string)
  default     = {}
}
