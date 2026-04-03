variable "name_prefix" {
  type = string
}

variable "cluster_name" {
  type = string
}

variable "db_identifier" {
  type = string
}

variable "api_cpu_threshold" {
  type    = number
  default = 80
}

variable "db_cpu_threshold" {
  type    = number
  default = 80
}

variable "alarm_actions" {
  type    = list(string)
  default = []
}

variable "tags" {
  type    = map(string)
  default = {}
}
