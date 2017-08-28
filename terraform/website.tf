resource "aws_s3_bucket" "default" {
  bucket = "charlieegan3-${var.project}-web"
  acl    = "public-read"

  force_destroy = true

  policy = <<EOF
{
        "Version": "2008-10-17",
        "Statement": [
                {
                        "Effect": "Allow",
                        "Principal": {
                                "AWS": "*"
                        },
                        "Action": "s3:GetObject",
                        "Resource": "arn:aws:s3:::charlieegan3-${var.project}-web/*"
                }
        ]
}
  EOF

  website {
    index_document = "index.html"
  }
}

resource "aws_s3_bucket_object" "index_page" {
  bucket       = "${aws_s3_bucket.default.id}"
  key          = "index.html"
  source       = "../web/index.html"
  etag         = "${md5(file("../web/index.html"))}"
  content_type = "text/html"
}

resource "aws_cloudfront_distribution" "website" {
  enabled = true
  aliases = ["${var.domain}"]

  origin {
    domain_name = "${aws_s3_bucket.default.website_endpoint}"
    origin_id   = "${aws_s3_bucket.default.id}"

    custom_origin_config {
      origin_protocol_policy = "http-only"
      http_port              = 80
      https_port             = 443
      origin_ssl_protocols   = ["TLSv1.2", "TLSv1.1", "TLSv1"]
    }
  }

  is_ipv6_enabled = true

  price_class = "PriceClass_100"

  "restrictions" {
    "geo_restriction" {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    acm_certificate_arn = "${data.aws_acm_certificate.default.arn}"

    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1"
  }

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "${aws_s3_bucket.default.id}"

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "redirect-to-https"

    min_ttl     = 0
    default_ttl = 3600
    max_ttl     = 18144000

    compress = true
  }
}
