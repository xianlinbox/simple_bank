terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    bucket = "mybucket"
    key    = "path/to/my/key"
    region = "${var.region}"
  }
}

# Configure the AWS Provider
provider "aws" {
  region = "${var.region}"
  profile = "tf-user"
}