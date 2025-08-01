// Package extip provides a Terraform provider for retrieving external IP addresses.
package extip

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{},

		DataSourcesMap: map[string]*schema.Resource{
			"extip": dataSource(),
		},

		ResourcesMap: map[string]*schema.Resource{},
	}
}
