package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBackendsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "trinogateway_backend" "this" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of clusters returned
					resource.TestCheckResourceAttr("data.trinogateway_backend.this", "clusters.#", "9"),
					// Verify the first cluster to ensure all attributes are set
					resource.TestCheckResourceAttr("data.trinogateway_backend.this", "clusters.0.description", ""),
					resource.TestCheckResourceAttr("data.trinogateway_backend.this", "clusters.0.id", "1"),
					resource.TestCheckResourceAttr("data.trinogateway_backend.this", "clusters.0.image", "/hashicorp.png"),
					resource.TestCheckResourceAttr("data.trinogateway_backend.this", "clusters.0.ingredients.#", "1"),
					resource.TestCheckResourceAttr("data.trinogateway_backend.this", "clusters.0.ingredients.0.id", "6"),
					resource.TestCheckResourceAttr("data.trinogateway_backend.this", "clusters.0.name", "HCP Aeropress"),
					resource.TestCheckResourceAttr("data.trinogateway_backend.this", "clusters.0.price", "200"),
					resource.TestCheckResourceAttr("data.trinogateway_backend.this", "clusters.0.teaser", "Automation in a cup"),
				),
			},
		},
	})
}
