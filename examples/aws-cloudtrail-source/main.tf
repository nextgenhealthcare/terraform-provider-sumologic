variable "auth_token" {}
variable "endpoint_url" {}
variable "sumologic_aws_external_id" {}

provider "sumologic" {
    auth_token = "${var.auth_token}"
    endpoint_url = "${var.endpoint_url}"
}

provider "aws" {
  region = "us-west-2"
}

resource "random_pet" "name" {}
