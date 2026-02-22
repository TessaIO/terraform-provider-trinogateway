// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/TessaIO/terraform-provider-trinogateway/internal/trinogateway"
)

func TestAccBackendsDataSource(t *testing.T) {
	testBackend := trinogateway.Backend{
		Name:         "cluster1",
		ProxyTo:      "localhost:8081",
		Active:       true,
		RoutingGroup: "adhoc",
		ExternalURL:  "trino.io",
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccBackendDataSource(testBackend),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of clusters returned
					resource.TestCheckResourceAttr("data.trinogateway_backends.all", "clusters.#", "1"),
					// Verify the first cluster to ensure all attributes are set
					resource.TestCheckResourceAttr("data.trinogateway_backends.all", "clusters.0.name", testBackend.Name),
					resource.TestCheckResourceAttr("data.trinogateway_backends.all", "clusters.0.routing_group", testBackend.RoutingGroup),
					resource.TestCheckResourceAttr("data.trinogateway_backends.all", "clusters.0.proxy_to", testBackend.ProxyTo),
					resource.TestCheckResourceAttr("data.trinogateway_backends.all", "clusters.0.active", fmt.Sprintf("%v", testBackend.Active)),
					resource.TestCheckResourceAttr("data.trinogateway_backends.all", "clusters.0.external_url", testBackend.ExternalURL),
				),
			},
		},
	})
}

func testAccBackendDataSource(backend trinogateway.Backend) string {
	return providerConfig + fmt.Sprintf(`
resource "trinogateway_backend" "this" {
  name          = "%s"
  routing_group = "%s"
  proxy_to      = "%s"
  active        = %v
  external_url  = "%s"
}

data "trinogateway_backends" "all" {
  depends_on = [trinogateway_backend.this]
}
`, backend.Name, backend.RoutingGroup, backend.ProxyTo, backend.Active, backend.ExternalURL)
}
