package sumologic

import (
	"github.com/brandonstevens/sumologic-sdk-go"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_token": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Authentication Token from Sumo Logic.",
				DefaultFunc: schema.EnvDefaultFunc("SUMOLOGIC_AUTH_TOKEN", nil),
			},
			"endpoint_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sumo Logic API Endpoint URL.",
				DefaultFunc: schema.EnvDefaultFunc("SUMOLOGIC_ENDPOINT_URL", "https://api.sumologic.com/api/v1/"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"sumologic_hosted_collector": resourceHostedCollector(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return sumologic.NewClient(
		d.Get("auth_token").(string),
		d.Get("endpoint_url").(string),
	)
}
