package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/nextgenhealthcare/terraform-provider-sumologic/sumologic"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: sumologic.Provider,
	})
}
