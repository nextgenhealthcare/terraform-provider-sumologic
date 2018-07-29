# Sumo Logic Provider Examples

This directory contains a set of examples of managing various Sumo Logic resources with Terraform. The examples each have their own README containing more details on what the example does.

To run the examples, you will need a Sumo Logic account. All the examples can be run on the free tier.

You will also need to generate a set of access keys and Base64 encode them. It's recommended that you create a new user with only Manage Collectors privileges. For more information, see [API Authentication](https://help.sumologic.com/APIs/General-API-Information/API-Authentication)

Additionally, you will need to know your API endpoint URL. Sumo Logic has several deployments that are assigned depending on the geographic location and the date an account is created. For more information, see [Sumo Logic Endpoints and Firewall Security](https://help.sumologic.com/APIs/General-API-Information/Sumo-Logic-Endpoints-and-Firewall-Security)

* [Creating an AWS CloudTrail Source in Sumo Logic](aws-cloudtrail-source/README.md)
* [Creating a Hosted Collector in Sumo Logic](hosted-collector/README.md)
