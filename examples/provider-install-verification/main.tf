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

data "trinogateway_cluster" "all" {}

output "clusters" {
  value = data.trinogateway_cluster.all
}
