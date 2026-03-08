terraform {
  required_providers {
    trinogateway = {
      source  = "tessaio/trinogateway"
      version = "0.1.0"
    }
  }
}

provider "trinogateway" {
  endpoint = "http://localhost:8080"
  username = "admin"
  password = "admin"
}

resource "trinogateway_backend" "this" {
  name          = "test-cluster2"
  routing_group = "test"
  proxy_to      = "http://localhost:8081"
  active        = false
  external_url  = "http://google.com"
}

data "trinogateway_backends" "all" {
  depends_on = [trinogateway_backend.this]
}

output "backends" {
  value = data.trinogateway_backends.all
}
