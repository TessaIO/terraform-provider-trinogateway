// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"regexp"
	"testing"


	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/TessaIO/terraform-provider-trinogateway/internal/trinogateway"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the TrinoGateway client is properly configured.
	// It is also possible to use the TRINOGATEWAY_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	providerConfig = `
provider "trinogateway" {
  endpoint = "http://localhost:8080"
  username = "admin"
  password = "admin"
}
`
)

// testAccProtoV6ProviderFactories is used to instantiate a provider during acceptance testing.
// The factory function is called for each Terraform CLI command to create a provider
// server that the CLI can connect to and interact with.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"trinogateway": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}

func TestAccProvider_configuredWithEnvVars(t *testing.T) {
	t.Setenv("TRINOGATEWAY_ENDPOINT", "http://localhost:8080")
	t.Setenv("TRINOGATEWAY_USERNAME", "admin")
	t.Setenv("TRINOGATEWAY_PASSWORD", "admin")

	testBackend := trinogateway.Backend{
		Name:         "cluster1",
		ProxyTo:      "localhost:8081",
		Active:       true,
		RoutingGroup: "adhoc",
		ExternalURL:  "trino.io",
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				// Test that we can read the backends data source without any explicit provider config
				Config: testAccProvider(testBackend),
				Check: resource.ComposeTestCheckFunc(
					// We don't know how many clusters there are, but we can check the attribute exists
					resource.TestCheckResourceAttrSet("data.trinogateway_backends.all", "clusters.#"),
				),
			},
		},
	})
}

func TestAccProvider_missingConfig(t *testing.T) {
	// Unset env vars to ensure a clean slate
	t.Setenv("TRINOGATEWAY_ENDPOINT", "")
	t.Setenv("TRINOGATEWAY_USERNAME", "")
	t.Setenv("TRINOGATEWAY_PASSWORD", "")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				Config: `
					provider "trinogateway" {}

					data "trinogateway_backends" "all" {}
				`,
				ExpectError: regexp.MustCompile("Error running pre-apply plan"),
			},
		},
	})
}

func TestProvider_Configure(t *testing.T) {
	t.Run("no_config", func(t *testing.T) {
		t.Setenv("TRINOGATEWAY_ENDPOINT", "")
		t.Setenv("TRINOGATEWAY_USERNAME", "")
		t.Setenv("TRINOGATEWAY_PASSWORD", "")
		var (
			p      = New("test")()
			resp   = &provider.ConfigureResponse{}
			reqCtx = context.Background()
			req    = provider.ConfigureRequest{
				Config: createProviderConfig(p, map[string]tftypes.Value{}),
			}
		)
		p.Configure(reqCtx, req, resp)
		if !resp.Diagnostics.HasError() {
			t.Fatal("expected an error, but got none")
		}
		if resp.Diagnostics.ErrorsCount() != 3 {
			t.Errorf("expected 3 errors, but got %d", resp.Diagnostics.ErrorsCount())
		}
	})

	t.Run("config_from_env", func(t *testing.T) {
		t.Setenv("TRINOGATEWAY_ENDPOINT", "http://localhost:8080")
		t.Setenv("TRINOGATEWAY_USERNAME", "user")
		t.Setenv("TRINOGATEWAY_PASSWORD", "pass")

		var (
			p      = New("test")()
			resp   = &provider.ConfigureResponse{}
			reqCtx = context.Background()
			req    = provider.ConfigureRequest{
				Config: createProviderConfig(p, map[string]tftypes.Value{}),
			}
		)

		p.Configure(reqCtx, req, resp)

		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
		if resp.DataSourceData == nil {
			t.Fatal("DataSourceData should be set")
		}
		if resp.ResourceData == nil {
			t.Fatal("ResourceData should be set")
		}
	})

	t.Run("config_from_provider_args", func(t *testing.T) {
		var (
			p      = New("test")()
			resp   = &provider.ConfigureResponse{}
			reqCtx = context.Background()
			req    = provider.ConfigureRequest{
				Config: createProviderConfig(p, map[string]tftypes.Value{
					"endpoint": tftypes.NewValue(tftypes.String, "http://localhost:8080"),
					"username": tftypes.NewValue(tftypes.String, "user"),
					"password": tftypes.NewValue(tftypes.String, "pass"),
				}),
			}
		)

		p.Configure(reqCtx, req, resp)

		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
		if resp.DataSourceData == nil {
			t.Fatal("DataSourceData should be set")
		}
		if resp.ResourceData == nil {
			t.Fatal("ResourceData should be set")
		}
	})

	t.Run("provider_args_precedence", func(t *testing.T) {
		t.Setenv("TRINOGATEWAY_ENDPOINT", "http://env-localhost:8080")
		t.Setenv("TRINOGATEWAY_USERNAME", "env-user")
		t.Setenv("TRINOGATEWAY_PASSWORD", "env-pass")

		var (
			p      = New("test")()
			resp   = &provider.ConfigureResponse{}
			reqCtx = context.Background()
			req    = provider.ConfigureRequest{
				Config: createProviderConfig(p, map[string]tftypes.Value{
					"endpoint": tftypes.NewValue(tftypes.String, "http://arg-localhost:8080"),
					"username": tftypes.NewValue(tftypes.String, "arg-user"),
					"password": tftypes.NewValue(tftypes.String, "arg-pass"),
				}),
			}
		)

		p.Configure(reqCtx, req, resp)

		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
		// NOTE: We can't easily assert the values used for configuration without mocking the client.
		// This test mainly ensures that the configuration is successful and doesn't throw errors.
	})

	t.Run("missing_username_password", func(t *testing.T) {
		t.Setenv("TRINOGATEWAY_USERNAME", "")
		t.Setenv("TRINOGATEWAY_PASSWORD", "")
		var (
			p      = New("test")()
			resp   = &provider.ConfigureResponse{}
			reqCtx = context.Background()
			req    = provider.ConfigureRequest{
				Config: createProviderConfig(p, map[string]tftypes.Value{
					"endpoint": tftypes.NewValue(tftypes.String, "http://localhost:8080"),
				}),
			}
		)
		p.Configure(reqCtx, req, resp)

		if !resp.Diagnostics.HasError() {
			t.Fatal("expected an error, but got none")
		}
		if resp.Diagnostics.ErrorsCount() != 2 {
			t.Errorf("expected 2 errors, but got %d", resp.Diagnostics.ErrorsCount())
		}
	})

	t.Run("missing_password", func(t *testing.T) {
		t.Setenv("TRINOGATEWAY_PASSWORD", "")
		var (
			p      = New("test")()
			resp   = &provider.ConfigureResponse{}
			reqCtx = context.Background()
			req    = provider.ConfigureRequest{
				Config: createProviderConfig(p, map[string]tftypes.Value{
					"endpoint": tftypes.NewValue(tftypes.String, "http://localhost:8080"),
					"username": tftypes.NewValue(tftypes.String, "user"),
				}),
			}
		)
		p.Configure(reqCtx, req, resp)

		if !resp.Diagnostics.HasError() {
			t.Fatal("expected an error, but got none")
		}
		if resp.Diagnostics.ErrorsCount() != 1 {
			t.Errorf("expected 1 error, but got %d", resp.Diagnostics.ErrorsCount())
		}
	})
}

// createProviderConfig is a helper function to create the provider configuration for testing.
func createProviderConfig(p provider.Provider, vals map[string]tftypes.Value) tfsdk.Config {
	schemaResp := &provider.SchemaResponse{}
	p.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)

	objVals := make(map[string]tftypes.Value)

	// Set explicit values provided by the test
	for k, v := range vals {
		objVals[k] = v
	}

	// Fill in missing values with nulls
	if _, ok := objVals["endpoint"]; !ok {
		objVals["endpoint"] = tftypes.NewValue(tftypes.String, nil)
	}
	if _, ok := objVals["username"]; !ok {
		objVals["username"] = tftypes.NewValue(tftypes.String, nil)
	}
	if _, ok := objVals["password"]; !ok {
		objVals["password"] = tftypes.NewValue(tftypes.String, nil)
	}

	return tfsdk.Config{
		Raw: tftypes.NewValue(tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"endpoint": tftypes.String,
				"username": tftypes.String,
				"password": tftypes.String,
			},
		}, objVals),
		Schema: schemaResp.Schema,
	}
}

func testAccProvider(backend trinogateway.Backend) string {
	return providerConfig + fmt.Sprintf(`
resource "trinogateway_backend" "this" {
  name          = "%s"
  routing_group = "%s"
  proxy_to      = "%s"
  active        = "%v"
  external_url  = "%s"
}

data "trinogateway_backends" "all" {
  depends_on = [trinogateway_backend.this]
}
`, backend.Name, backend.RoutingGroup, backend.ProxyTo, backend.Active, backend.ExternalURL)
}
