variable "auth_token" {}
variable "endpoint_url" {}

provider "sumologic" {
    auth_token = "${var.auth_token}"
    endpoint_url = "${var.endpoint_url}"
}

resource "random_pet" "name" {}

resource "sumologic_hosted_collector" "example" {
    name = "example-${random_pet.name.id}"
}
