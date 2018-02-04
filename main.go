package main

import (
	"github.com/brandonstevens/terraform-provider-sumologic/sumologic"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: sumologic.Provider,
	})
}
