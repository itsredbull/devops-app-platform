resource "aws_db_subnet_group" "this" {
  name       = "${var.name_prefix}-db-subnet-group"
  subnet_ids = var.subnet_ids

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-db-subnet-group"
  })
}

resource "aws_security_group" "db" {
  name        = "${var.name_prefix}-db-sg"
  description = "Postgres access for uptime platform"
  vpc_id      = var.vpc_id

  ingress {
    description = "Postgres from VPC"
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [var.vpc_cidr]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-db-sg"
  })
}

resource "aws_db_instance" "this" {
  identifier                 = "${var.name_prefix}-postgres"
  allocated_storage          = var.allocated_storage
  engine                     = "postgres"
  engine_version             = var.engine_version
  instance_class             = var.instance_class
  db_name                    = var.db_name
  username                   = var.db_username
  password                   = var.db_password
  db_subnet_group_name       = aws_db_subnet_group.this.name
  vpc_security_group_ids     = [aws_security_group.db.id]
  backup_retention_period    = var.backup_retention_days
  skip_final_snapshot        = var.skip_final_snapshot
  storage_encrypted          = true
  publicly_accessible        = false
  deletion_protection        = var.deletion_protection
  auto_minor_version_upgrade = true

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-postgres"
  })
}
