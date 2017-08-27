variable project {
  default = "rssmerge"
}

variable domain {
  default = "rssmerge.charlieegan3.com"
}

data "aws_caller_identity" "current" {}

variable "region" {
  default = "us-east-1"
}

terraform {
  backend "s3" {
    bucket = "charlieegan3-www-terraform-state"
    region = "us-east-1"
    key    = "rssmerge.tfstate"
  }
}
