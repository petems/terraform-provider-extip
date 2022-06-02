package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/petems/terraform-provider-extip/extip"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: extip.Provider})
}
