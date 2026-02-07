terraform {
  required_providers {
    trino-gateway = {
      source  = "trino-gateway"
      version = "1.0.0"
    }
  }
}

provider "trino-gateway" {
  endpoint = "localhost:8080"
}

# data "trino_gateway_cluster" "main" {}
