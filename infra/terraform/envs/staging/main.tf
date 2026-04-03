module "network" {
  source = "../../modules/network"

  name_prefix          = var.name_prefix
  vpc_cidr             = var.vpc_cidr
  public_subnet_cidrs  = var.public_subnet_cidrs
  private_subnet_cidrs = var.private_subnet_cidrs
  availability_zones   = var.availability_zones
  enable_nat_gateway   = var.enable_nat_gateway
  tags                 = var.tags
}

module "kubernetes" {
  source = "../../modules/kubernetes"

  name_prefix        = var.name_prefix
  cluster_role_arn   = var.eks_cluster_role_arn
  node_role_arn      = var.eks_node_role_arn
  subnet_ids         = module.network.private_subnet_ids
  kubernetes_version = "1.30"
  desired_nodes      = 3
  min_nodes          = 2
  max_nodes          = 5
  node_instance_types = [
    "t3.large"
  ]
  tags = var.tags
}

module "database" {
  source = "../../modules/database"

  name_prefix            = var.name_prefix
  vpc_id                 = module.network.vpc_id
  vpc_cidr               = module.network.vpc_cidr
  subnet_ids             = module.network.private_subnet_ids
  db_name                = var.db_name
  db_username            = var.db_username
  db_password            = var.db_password
  instance_class         = "db.t4g.small"
  allocated_storage      = 50
  backup_retention_days  = 7
  skip_final_snapshot    = false
  deletion_protection    = true
  tags                   = var.tags
}

module "monitoring" {
  source = "../../modules/monitoring"

  name_prefix       = var.name_prefix
  cluster_name      = module.kubernetes.cluster_name
  db_identifier     = module.database.db_identifier
  api_cpu_threshold = 75
  db_cpu_threshold  = 75
  alarm_actions     = var.alarm_actions
  tags              = var.tags
}
