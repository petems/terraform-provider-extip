# terraform-provider-extip

Terraform provider for getting your current external IP as a data source.

## Requirements
-	[Terraform](https://www.terraform.io/downloads.html) 0.11.x
-	[Go](https://golang.org/doc/install) 1.10 (to build the provider plugin)

## Installing the Provider
Follow the instructions to [install it as a plugin](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin). After placing it into your plugins directory, run `terraform init` to initialize it.

## Usage
```
data "extip" "external_ip" {
}

output "external_ip" {
  value = "${data.extip.external_ip.ipaddress}"
}

```

Gives the result:
```
data.extip.external_ip: Refreshing state...

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

external_ip = 238.209.109.16
```

You can also specify what resolver you want to use to get the URL:

```
data "extip" "external_ip_from_aws" {
  resolver = "https://checkip.amazonaws.com/"
}

output "external_ip_from_aws" {
  value = "${data.extip.external_ip_from_aws.ipaddress}"
}
```


Examples are under [/examples](/examples).

## Building the Provider
Clone and build the repository

```sh
go get github.com/petems/terraform-provider-extip
make build
```

Symlink the binary to your terraform plugins directory:

```sh
ln -s $GOPATH/bin/terraform-provider-extip ~/.terraform.d/plugins/
```

## Updating the Provider

```sh
go get -u github.com/petems/terraform-provider-extip
make build
```

## TODO

* Add configuration of the consensus timing (ie. how long it will wait to resolve)
* Add option of getting ipv6 or ipv4 ipaddress

## Contributing
* Install project dependencies: `go get github.com/kardianos/govendor`
* Run tests: `make test`
* Build the binary: `make build`
