package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/viktorradnai/terraform-provider-bcrypt/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
