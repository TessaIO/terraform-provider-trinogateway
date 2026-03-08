# Trino Gateway Terraform Provider

A Terraform provider for managing [Trino Gateway](https://github.com/trinodb/trino-gateway).

This provider allows you to manage Trino Gateway backend configurations using Terraform.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24 (for development)

## Usage Example

Here is a complete example of how to configure the provider and manage a Trino Gateway backend.

```hcl
terraform {
  required_providers {
    trinogateway = {
      source  = "tessaio/trinogateway"
      version = "0.1.0"
    }
  }
}

provider "trinogateway" {
  # It is recommended to use environment variables for credentials
  # endpoint = "http://localhost:8080"
  # username = "admin"
  # password = "admin"
}

resource "trinogateway_backend" "cluster1" {
  name          = "cluster1"
  proxy_to      = "http://trino-coordinator-1:8080"
  active        = true
  routing_group = "adhoc"
  external_url  = "http://trino.example.com/cluster1"
}

data "trinogateway_backends" "all" {
  depends_on = [trinogateway_backend.cluster1]
}

output "all_backends" {
  value = data.trinogateway_backends.all.clusters
}
```

## Provider Configuration

The Trino Gateway provider is configured using the following arguments:

### Argument Reference

- `endpoint` - (Required) The endpoint of the Trino Gateway API. Can also be set with the `TRINOGATEWAY_ENDPOINT` environment variable.
- `username` - (Required) The username for authenticating with the Trino Gateway API. Can also be set with the `TRINOGATEWAY_USERNAME` environment variable.
- `password` - (Required) The password for authenticating with the Trino Gateway API. Can also be set with the `TRINOGATEWAY_PASSWORD` environment variable.

## Resources

### trinogateway_backend

Manages a backend configuration in Trino Gateway.

#### Example Usage

```hcl
resource "trinogateway_backend" "etl_cluster" {
  name          = "etl_cluster"
  proxy_to      = "http://trino-coordinator-etl:8080"
  active        = true
  routing_group = "etl"
  external_url  = "http://trino.example.com/etl"
}
```

#### Argument Reference

- `name` - (Required) The unique name of the backend.
- `proxy_to` - (Required) The URL that the gateway will proxy requests to for this backend.
- `active` - (Required) A boolean flag indicating whether the backend is active.
- `routing_group` - (Required) The routing group to which this backend belongs.
- `external_url` - (Optional) The external URL that can be used to access this backend.

#### Attribute Reference

In addition to all the arguments above, the following attributes are exported:

- `name` - The name of the backend.
- `proxy_to` - The proxy URL for the backend.
- `active` - Whether the backend is active.
- `routing_group` - The routing group for the backend.
- `external_url` - The external URL for the backend.

## Data Sources

### trinogateway_backends

Provides a list of all backend configurations in Trino Gateway.

#### Example Usage

```hcl
data "trinogateway_backends" "all" {}

output "backend_names" {
  value = [for backend in data.trinogateway_backends.all.clusters : backend.name]
}
```

#### Attribute Reference

- `clusters` - A list of backend objects. Each object has the following attributes:
    - `name` - The name of the backend.
    - `proxy_to` - The proxy URL for the backend.
    - `active` - Whether the backend is active.
    - `routing_group` - The routing group for the backend.
    - `external_url` - The external URL for the backend.

## Roadmap

The current version of the provider delivers foundational support for managing Trino Gateway backends.

-   **`trinogateway_backend` Resource**: Full lifecycle management (create, read, update, delete, import) for backend configurations.
-   **`trinogateway_backends` Data Source**: A data source to read all backend configurations.

Future development will focus on expanding coverage of the Trino Gateway API. The immediate priority is to add support for:

-   **Routing Groups**: A new `trinogateway_routing_group` resource to manage routing groups.

We also plan to investigate and add support for other API resources as they are prioritized by the community.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine.

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources.

```shell
make testacc
```
