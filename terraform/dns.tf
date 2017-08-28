resource "aws_route53_record" "api" {
  zone_id = "${var.zone_id}"

  name = "${aws_api_gateway_domain_name.default.domain_name}"
  type = "A"

  alias {
    name                   = "${aws_api_gateway_domain_name.default.cloudfront_domain_name}"
    zone_id                = "${aws_api_gateway_domain_name.default.cloudfront_zone_id}"
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "website" {
  zone_id = "${var.zone_id}"

  name = "${var.domain}"
  type = "A"

  alias {
    name                   = "${aws_cloudfront_distribution.website.domain_name}"
    zone_id                = "${aws_cloudfront_distribution.website.hosted_zone_id}"
    evaluate_target_health = true
  }
}
