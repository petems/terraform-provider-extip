package extip

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{},

		DataSourcesMap: map[string]*schema.Resource{
			"extip": dataSource(),
		},

		ResourcesMap: map[string]*schema.Resource{},
	}
}
