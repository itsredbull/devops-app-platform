resource "aws_eks_cluster" "this" {
  name     = "${var.name_prefix}-eks"
  role_arn = var.cluster_role_arn
  version  = var.kubernetes_version

  vpc_config {
    subnet_ids              = var.subnet_ids
    endpoint_public_access  = true
    endpoint_private_access = true
  }

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-eks"
  })
}

resource "aws_eks_node_group" "default" {
  cluster_name    = aws_eks_cluster.this.name
  node_group_name = "${var.name_prefix}-nodes"
  node_role_arn   = var.node_role_arn
  subnet_ids      = var.subnet_ids

  scaling_config {
    desired_size = var.desired_nodes
    min_size     = var.min_nodes
    max_size     = var.max_nodes
  }

  instance_types = var.node_instance_types

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-nodes"
  })
}
