locals {
    tags = {
    owner       = "xianlinbox"
    project    = "simple-bank"
  }
}
resource "aws_ecr_repository" "simple_bank_ecr_repo" {
  name                 = "simple-bank-repo" 
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
  tags = local.tags
}
