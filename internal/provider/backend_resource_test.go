// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/TessaIO/terraform-provider-trinogateway/internal/client"
	"github.com/TessaIO/terraform-provider-trinogateway/internal/trinogateway"
)

func TestAccBackendResource(t *testing.T) {
	rName := "cluster1"
	resourceName := "trinogateway_backend." + rName

	// Initial backend configuration
	testBackendCreate := trinogateway.Backend{
		Name:         rName,
		ProxyTo:      "http://localhost:8081",
		Active:       true,
		RoutingGroup: "adhoc",
		ExternalURL:  "http://trino.example.com",
	}

	// Updated backend configuration
	testBackendUpdate := trinogateway.Backend{
		Name:         rName,
		ProxyTo:      "http://localhost:8082",
		Active:       true,
		RoutingGroup: "etl",
		ExternalURL:  "http://trino-updated.example.com",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckBackendResourceDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create and Read testing
			{
				Config: testAccBackendResourceConfig(rName, testBackendCreate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", testBackendCreate.Name),
					resource.TestCheckResourceAttr(resourceName, "proxy_to", testBackendCreate.ProxyTo),
					resource.TestCheckResourceAttr(resourceName, "active", fmt.Sprintf("%v", testBackendCreate.Active)),
					resource.TestCheckResourceAttr(resourceName, "routing_group", testBackendCreate.RoutingGroup),
					resource.TestCheckResourceAttr(resourceName, "external_url", testBackendCreate.ExternalURL),
				),
			},
			// Step 2: Update and Read testing
			{
				Config: testAccBackendResourceConfig(rName, testBackendUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", testBackendUpdate.Name),
					resource.TestCheckResourceAttr(resourceName, "proxy_to", testBackendUpdate.ProxyTo),
					resource.TestCheckResourceAttr(resourceName, "active", fmt.Sprintf("%v", testBackendUpdate.Active)),
					resource.TestCheckResourceAttr(resourceName, "routing_group", testBackendUpdate.RoutingGroup),
					resource.TestCheckResourceAttr(resourceName, "external_url", testBackendUpdate.ExternalURL),
				),
			},
		},
	})
}

func testAccBackendResourceConfig(rName string, backend trinogateway.Backend) string {
	return providerConfig + fmt.Sprintf(`
resource "trinogateway_backend" "%s" {
  name          = "%s"
  proxy_to      = "%s"
  active        = %v
  routing_group = "%s"
  external_url  = "%s"
}
`, rName, backend.Name, backend.ProxyTo, backend.Active, backend.RoutingGroup, backend.ExternalURL)
}

func testAccCheckBackendResourceDestroy(s *terraform.State) error {
	// TODO: Make these shared across all tests
	endpoint := "http://localhost:8080"
	username := "admin"
	password := "admin"

	httpClient, err := client.NewClient(endpoint, client.WithAuth(username, password))
	if err != nil {
		return err
	}
	trinoGatewayService := trinogateway.NewTrinoGatewayService(httpClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "trinogateway_backend" {
			continue
		}

		_, err := trinoGatewayService.GetBackend(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Backend '%s' still exists.", rs.Primary.ID)
		}
	}

	return nil
}
