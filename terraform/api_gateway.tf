output "api_endpoint" {
  value = "${aws_api_gateway_deployment.default.invoke_url}/${var.api_gateway_path_part}"
}

variable "api_gateway_stage" {
  default = "production"
}

variable "api_gateway_path_part" {
  default = "build"
}

resource "aws_api_gateway_rest_api" "default" {
  name        = "${var.project}"
  description = "API for the ${var.project} lambda function"
}

resource "aws_api_gateway_deployment" "default" {
  depends_on  = ["aws_api_gateway_integration.default"]
  rest_api_id = "${aws_api_gateway_rest_api.default.id}"
  stage_name  = "${var.api_gateway_stage}"

  // https://github.com/hashicorp/terraform/issues/6613
  stage_description = "${md5(file("api_gateway.tf"))}"
}

resource "aws_api_gateway_domain_name" "default" {
  domain_name     = "api-${var.domain}"
  certificate_arn = "${data.aws_acm_certificate.default.arn}"
}

resource "aws_api_gateway_base_path_mapping" "default" {
  api_id      = "${aws_api_gateway_rest_api.default.id}"
  stage_name  = "${aws_api_gateway_deployment.default.stage_name}"
  domain_name = "${aws_api_gateway_domain_name.default.domain_name}"
}

resource "aws_api_gateway_method_settings" "default" {
  rest_api_id = "${aws_api_gateway_rest_api.default.id}"
  stage_name  = "${var.api_gateway_stage}"
  method_path = "${aws_api_gateway_resource.default.path_part}/*"

  settings {
    metrics_enabled = true
    logging_level   = "INFO"
  }

  depends_on = ["aws_api_gateway_deployment.default"]
}

resource "aws_api_gateway_resource" "default" {
  rest_api_id = "${aws_api_gateway_rest_api.default.id}"
  parent_id   = "${aws_api_gateway_rest_api.default.root_resource_id}"
  path_part   = "${var.api_gateway_path_part}"
}

resource "aws_api_gateway_method" "default" {
  rest_api_id   = "${aws_api_gateway_rest_api.default.id}"
  resource_id   = "${aws_api_gateway_resource.default.id}"
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "default" {
  rest_api_id             = "${aws_api_gateway_rest_api.default.id}"
  resource_id             = "${aws_api_gateway_resource.default.id}"
  http_method             = "${aws_api_gateway_method.default.http_method}"
  type                    = "AWS_PROXY"
  uri                     = "arn:aws:apigateway:${var.region}:lambda:path/2015-03-31/functions/${aws_lambda_function.lambda.arn}/invocations"
  integration_http_method = "POST"
}

resource "aws_api_gateway_method_response" "response_method" {
  rest_api_id = "${aws_api_gateway_rest_api.default.id}"
  resource_id = "${aws_api_gateway_resource.default.id}"
  http_method = "${aws_api_gateway_integration.default.http_method}"
  status_code = "200"
}

resource "aws_api_gateway_integration_response" "response_method_integration" {
  rest_api_id = "${aws_api_gateway_rest_api.default.id}"
  resource_id = "${aws_api_gateway_resource.default.id}"
  http_method = "${aws_api_gateway_method_response.response_method.http_method}"
  status_code = "${aws_api_gateway_method_response.response_method.status_code}"
}
