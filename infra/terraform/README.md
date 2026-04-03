# Terraform Infrastructure

This directory contains reusable Terraform modules and environment compositions for `dev` and `staging`.

## Layout

- `global/`: bootstrap remote state infrastructure (S3 + DynamoDB)
- `modules/network`: VPC, subnets, IGW, and public route table
- `modules/kubernetes`: EKS cluster and managed node group
- `modules/database`: PostgreSQL RDS and DB subnet/security group
- `modules/monitoring`: CloudWatch alarms for EKS and RDS
- `envs/dev`: development environment stack
- `envs/staging`: staging environment stack

## 1) Bootstrap remote state

```bash
cd infra/terraform/global
cp terraform.tfvars.example terraform.tfvars
terraform init
terraform apply
```

## 2) Deploy an environment

```bash
cd infra/terraform/envs/dev
cp backend.hcl.example backend.hcl
cp terraform.tfvars.example terraform.tfvars
terraform init -backend-config=backend.hcl
terraform plan
terraform apply
```

For staging, use `infra/terraform/envs/staging` and the staging `backend.hcl` key.

## Notes

- IAM roles for EKS cluster and nodes must exist before apply.
- Keep real `terraform.tfvars` and `backend.hcl` files out of git.
- Use different state keys per environment to avoid collisions.
