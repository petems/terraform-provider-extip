terraform {
  required_providers {
    extip = {
      source = "petems/extip"
      version = "0.1.0"
    }
  }
}

data "extip" "external_ip" {
}

data "extip" "external_ip_from_aws" {
  resolver = "https://checkip.amazonaws.com/"
}

output "external_ip" {
  value = data.extip.external_ip.ipaddress
}

output "external_ip_from_aws" {
  value = data.extip.external_ip_from_aws.ipaddress
}
