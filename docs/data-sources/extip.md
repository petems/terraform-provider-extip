---
layout: ""
page_title: "Provider: extip"
description: |-
  The Extip provider provides resources to return an external IP from a resolver.
---

# Data Source `extip`

The `extip` data source returns an IP address from an external resolver

## Example Usage

```terraform
data "extip" "check_external_ip" {
  resolver       = "https://checkip.example.com/"
  client_timeout = 500
}

output "check_external_ip" {
  value = data.extip.check_external_ip.ipaddress
}
```

## Schema

### Optional

- `resolver` (String, Optional) The address to use as a resolver. Defaults to `https://checkip.amazonaws.com/`
- `client_timeout` (Integer, Optional) The time to wait for a response in ms. If not set, defaults to `1000` (1 second). Setting to `0` means infinite (no timeout)
- `validate_ip` - (Boolean, Optional) Validates if the returned response is a valid ip address

### Read-only

- `ipaddress` The IP address returned from the resolver
