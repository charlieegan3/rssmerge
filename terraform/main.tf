variable project {
  default = "rssmerge"
}

variable domain {
  default = "rssmerge.charlieegan3.com"
}

variable "zone_id" {}

variable "region" {
  default = "us-east-1"
}

data "aws_caller_identity" "current" {}

data "aws_acm_certificate" "default" {
  domain   = "charlieegan3.com"
  statuses = ["ISSUED"]
}

terraform {
  backend "s3" {
    bucket = "charlieegan3-www-terraform-state"
    region = "us-east-1"
    key    = "rssmerge.tfstate"
  }
}
