data "extip" "external_ip" {
}

data "extip" "external_ip_from_aws" {
  resolver = "https://checkip.amazonaws.com/"
}

output "external_ip" {
  value = "${data.extip.external_ip.ipaddress}"
}

output "external_ip_from_aws" {
  value = "${data.extip.external_ip_from_aws.ipaddress}"
}
