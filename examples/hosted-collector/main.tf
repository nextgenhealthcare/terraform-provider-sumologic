variable "auth_token" {}
variable "endpoint_url" {}
variable "collector_name" {}

provider "sumologic" {
    auth_token = "${var.auth_token}"
    endpoint_url = "${var.endpoint_url}"
}

resource "sumologic_hosted_collector" "example" {
    name = "${var.collector_name}"
}