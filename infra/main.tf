locals {
  tags = {
    owner   = "xianlinbox"
    project = "simple-bank"
  }
}

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    bucket = "xianlinbox-simple-bank"
    key    = "tf-state/terraform.tfstate"
    region = "us-east-1"
  }
}

# Configure the AWS Provider
provider "aws" {
  region = "us-east-2"
}

resource "aws_ecr_repository" "simple_bank_ecr_repo" {
  name                 = "simple-bank-repo"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
  tags = local.tags
}


module "secrets" {
  source      = "./modules/secrets"
  DB_PASSWORD = var.DB_PASSWORD
  tags        = local.tags
}

resource "aws_db_instance" "simple_bank_db" {
  allocated_storage   = 10
  db_name             = "simple-bank"
  engine              = "postgres"
  engine_version      = "16.0"
  instance_class      = "db.t3.micro"
  username            = "root"
  password            = module.secrets.db_password_secret
  skip_final_snapshot = true

  depends_on = [
    module.secrets
  ]
}