---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "extip Data Source - terraform-provider-extip"
subcategory: ""
description: |-
  
---

# extip (Data Source)

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `client_timeout` (Number) The time to wait for a response in ms
If not set, defaults to 1000 (1 second). Setting to 0 means infinite (no timeout)
- `resolver` (String) The URL to use to resolve the external IP address
If not set, defaults to <https://checkip.amazonaws.com/>
- `validate_ip` (Boolean) Validate if the returned response is a valid ip address

### Read-Only

- `id` (String) The ID of this resource.
- `ipaddress` (String)
