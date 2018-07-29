## Creating a Hosted Collector in Sumo Logic

This example provides sample configuration for creating a hosted collector.

Initialize by running `terraform init`.

Once ready, run `terraform plan -out example.plan` to review.

You will be prompted to provide input for the following variables:

* auth_token: Base64 encoding of `<accessId>:<accessKey>`. For more information, see [API Authentication](https://help.sumologic.com/APIs/General-API-Information/API-Authentication)
* endpoint_url: Sumo Logic has several deployments that are assigned depending on the geographic location and the date an account is created. For more information, see [Sumo Logic Endpoints and Firewall Security](https://help.sumologic.com/APIs/General-API-Information/Sumo-Logic-Endpoints-and-Firewall-Security)

Once satisfied with plan, run `terraform apply example.plan`  
