resource "trinogateway_backend" "this" {
  name          = "test-cluster2"
  routing_group = "test"
  proxy_to      = "http://localhost:8081"
  active        = false
  external_url  = "http://google.com"
}
