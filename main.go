// Package main provides the entry point for the Terraform extip provider.
package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/petems/terraform-provider-extip/extip"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: extip.Provider,
	}

	if debug {
		opts.Debug = true
		opts.ProviderAddr = "registry.terraform.io/petems/extip"
	}

	plugin.Serve(opts)
}
