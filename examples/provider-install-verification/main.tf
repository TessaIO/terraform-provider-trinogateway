terraform {
  required_providers {
    trinogateway = {
      source  = "trinogateway"
      version = "1.0.0"
    }
  }
}

provider "trinogateway" {
  endpoint = "http://localhost:8080"
  username = "admin"
  password = "admin"
}

resource "trinogateway_cluster" "this" {
  name          = "test-cluster"
  routing_group = "test"
  proxy_to      = "http://localhost:8081"
  active        = "false"
  external_url  = "http://google.com"
}

data "trinogateway_cluster" "all" {}

output "clusters" {
  value = data.trinogateway_cluster.all
}
