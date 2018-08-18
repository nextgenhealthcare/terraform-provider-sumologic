resource "random_pet" "name" {}

resource "sumologic_hosted_collector" "example" {
    name = "example-${random_pet.name.id}"
}
