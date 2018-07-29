## Creating an AWS CloudTrail Source in Sumo Logic

This example provides sample configuration for creating an AWS CloudTrail source.

To run this example, in addition to a Sumo Logic account, you'll also need an AWS Account and credentials. The example does not specify AWS credential configuration and will either use the default profile or environment variables.

Initialize by running `terraform init`.

Once ready run `terraform plan -out example.plan` to review.

You will be prompted to provide input for the following variables:

* auth_token: Base64 encoding of `<accessId>:<accessKey>`. For more information, see [API Authentication](https://help.sumologic.com/APIs/General-API-Information/API-Authentication)
* endpoint_url: Sumo Logic has several deployments that are assigned depending on the geographic location and the date an account is created. For more information, see [Sumo Logic Endpoints and Firewall Security](https://help.sumologic.com/APIs/General-API-Information/Sumo-Logic-Endpoints-and-Firewall-Security)
* sumologic_aws_external_id: The External ID is unique to your Sumo account and needs to be in the specified format. For more information, see [Grant Access to an AWS Product](https://help.sumologic.com/Send-Data/Sources/02Sources-for-Hosted-Collectors/Amazon_Web_Services/Grant_Access_to_an_AWS_Product)

Once satisfied with plan, run `terraform apply example.plan`  
