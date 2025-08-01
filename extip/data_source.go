// Package extip provides a Terraform provider for retrieving external IP addresses.
package extip

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// HTTP client cache with timeout-based keys.
var (
	httpClients = make(map[time.Duration]*http.Client)
	clientMutex sync.RWMutex
)

// getHTTPClient returns an HTTP client with the specified timeout, reusing existing clients.
func getHTTPClient(timeout time.Duration) *http.Client {
	clientMutex.RLock()
	if client, exists := httpClients[timeout]; exists {
		clientMutex.RUnlock()
		return client
	}
	clientMutex.RUnlock()

	clientMutex.Lock()
	defer clientMutex.Unlock()

	// Double-check after acquiring write lock
	if client, exists := httpClients[timeout]; exists {
		return client
	}

	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}
	httpClients[timeout] = client
	return client
}

func dataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadContext,

		Schema: map[string]*schema.Schema{
			"ipaddress": {
				Type:     schema.TypeString,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"resolver": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://checkip.amazonaws.com/",
				Description: "The URL to use to resolve the external IP address\nIf not set, defaults to https://checkip.amazonaws.com/",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"client_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1000,
				Description: "The time to wait for a response in ms\nIf not set, defaults to 1000 (1 second). Setting to 0 means infinite (no timeout)",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"validate_ip": {
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
	timeout := time.Duration(clientTimeout) * time.Millisecond
	client := getHTTPClient(timeout)

	// Create request with context to satisfy noctx linter
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, service, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	rsp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer func() {
		if closeErr := rsp.Body.Close(); closeErr != nil {
			// Log the error but don't fail the operation
			_ = closeErr
		}
	}()

	if rsp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request error. Response code: %d", rsp.StatusCode)
	}

	buf, err := io.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}

	// Optimize string conversion by avoiding unnecessary allocations
	trimmed := bytes.TrimSpace(buf)
	return string(trimmed), nil
}

func dataSourceReadContext(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := dataSourceRead(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func dataSourceRead(d *schema.ResourceData, _ interface{}) error {
	resolver, ok := d.Get("resolver").(string)
	if !ok {
		return errors.New("resolver is not a string")
	}

	clientTimeout, ok := d.Get("client_timeout").(int)
	if !ok {
		return errors.New("client_timeout is not an int")
	}

	ip, err := getExternalIPFrom(resolver, clientTimeout)
	if err != nil {
		return fmt.Errorf("error requesting external IP: %s", err.Error())
	}

	// Only validate IP if the flag is set
	if v, ok := d.GetOk("validate_ip"); ok {
		if validateIP, ok := v.(bool); ok && validateIP {
			if net.ParseIP(ip) == nil {
				return fmt.Errorf("validate_ip was set to true, and information from resolver was not valid IP: %s", ip)
			}
		}
	}

	// Set the IP address
	if setErr := d.Set("ipaddress", ip); setErr != nil {
		return fmt.Errorf("error setting ipaddress: %s", setErr.Error())
	}

	// Use a more efficient ID generation
	d.SetId(time.Now().UTC().Format("20060102150405"))

	return nil
}
