package extip

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSource() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRead,

		Schema: map[string]*schema.Schema{
			"ipaddress": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"resolver": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://checkip.amazonaws.com/",
				Description: "The URL to use to resolve the external IP address\nIf not set, defaults to https://checkip.amazonaws.com/",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func getExternalIPFrom(service string) (string, error) {
	rsp, err := http.Get(service)
	if err != nil {
		return "", err
	}

	defer rsp.Body.Close()

	if rsp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP request error. Response code: %d", rsp.StatusCode)
	}

	buf, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}

	return string(bytes.TrimSpace(buf)), nil
}

func dataSourceRead(d *schema.ResourceData, meta interface{}) error {

	resolver := d.Get("resolver").(string)

	ip, err := getExternalIPFrom(resolver)

	if err == nil {
		d.Set("ipaddress", string(ip))
		d.SetId(time.Now().UTC().String())
	} else {
		return fmt.Errorf("Error requesting external IP: %d", err)
	}

	return nil

}
