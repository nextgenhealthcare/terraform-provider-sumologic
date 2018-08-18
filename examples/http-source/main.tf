resource "random_pet" "name" {}

resource "sumologic_hosted_collector" "example" {
    name = "example-${random_pet.name.id}"
}

resource "sumologic_http_source" "example" {
    name = "example"
    collector_id         = "${sumologic_hosted_collector.example.id}"
    source_type = "HTTP"
}

output "sumologic_http_source_url" {
    value = "${sumologic_http_source.example.url}"
}
