variable "sumologic_aws_external_id" {}

provider "aws" {
  region = "us-west-2"
}

resource "random_pet" "name" {}
