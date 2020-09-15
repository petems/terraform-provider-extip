package extip

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"client_timeout": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1000,
				Description: "The time to wait for a response in ms\nIf not set, defaults to 1000 (1 second). Setting to 0 means infinite (no timeout)",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"validate_ip": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Validate if the returned response is a valid ip address",
				Elem: &schema.Schema{
					Type: schema.TypeBool,
				},
			},
		},
	}
}

func getExternalIPFrom(service string, clientTimeout int) (string, error) {

	var netClient = &http.Client{
		Timeout: time.Duration(clientTimeout) * time.Millisecond,
	}

	rsp, err := netClient.Get(service)
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

	clientTimeout := d.Get("client_timeout").(int)

	ip, err := getExternalIPFrom(resolver, clientTimeout)

	if v, ok := d.GetOkExists("validate_ip"); ok {
		if v.(bool) {
			ipParse := net.ParseIP(ip)
			if ipParse == nil {
				return fmt.Errorf("validate_ip was set to true, and information from resolver was not valid IP: %s", ip)
			}
		}
	}

	if err == nil {
		d.Set("ipaddress", string(ip))
		d.SetId(time.Now().UTC().String())
	} else {
		return fmt.Errorf("Error requesting external IP: %s", err.Error())
	}

	return nil

}
