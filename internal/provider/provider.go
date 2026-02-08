// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/TessaIO/terraform-provider-trinogateway/internal/client"
	"github.com/TessaIO/terraform-provider-trinogateway/internal/trinogateway"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &TrinoGatewayProvider{}

// TrinoGatewayProvider defines the provider implementation.
type TrinoGatewayProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// TrinoGatewayProviderModel describes the provider data model.
type TrinoGatewayProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (p *TrinoGatewayProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "trinogateway"
	resp.Version = p.version
}

func (p *TrinoGatewayProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Endpoint of the Trino gateway",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "username to access Trino gateway",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "password to access Trino gateway",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *TrinoGatewayProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Trino-Gateway client")
	var config TrinoGatewayProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown Trino Gateway API Endpoint",
			"The provider cannot create the Trino-Gateway API client as there is an unknown configuration value for the Trino Gateway endpoint. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the TRINOGATEWAY_ENDPOINT environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Trino Gateway API Username",
			"The provider cannot create the Trino-Gateway API client as there is an unknown configuration value for the Trino Gateway username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the TRINOGATEWAY_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown Trino Gateway API Password",
			"The provider cannot create the Trino-Gateway API client as there is an unknown configuration value for the Trino Gateway password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the TRINOGATEWAY_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	endpoint := os.Getenv("TRINOGATEWAY_ENDPOINT")
	username := os.Getenv("TRINOGATEWAY_USERNAME")
	password := os.Getenv("TRINOGATEWAY_PASSWORD")

	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing Trino Gateway API Host",
			"The provider cannot create the Trino Gateway API client as there is a missing or empty value for the Trino Gateway API endpoint. "+
				"Set the host value in the configuration or use the TRINOGATEWAY_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Trino Gateway API Username",
			"The provider cannot create the Trino Gateway API client as there is a missing or empty value for the Trino Gateway API username. "+
				"Set the username value in the configuration or use the TRINOGATEWAY_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Trino Gateway API Password",
			"The provider cannot create the Trino Gateway API client as there is a missing or empty value for the Trino Gateway API password. "+
				"Set the password value in the configuration or use the TRINOGATEWAY_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "trinogateway_endpoint", endpoint)
	ctx = tflog.SetField(ctx, "trinogateway_username", username)
	ctx = tflog.SetField(ctx, "trinogateway_password", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "trinogateway_password")

	tflog.Debug(ctx, "Creating trinogateway client")

	httpClient, err := client.NewClient(endpoint, client.WithAuth(username, password))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create HTTP Client",
			"An unexpected error occurred when creating the HTTP client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"HTTP Client Error: "+err.Error(),
		)

		return
	}

	trinoGatewayService := trinogateway.NewTrinoGatewayService(httpClient)
	resp.ResourceData = trinoGatewayService
	resp.DataSourceData = trinoGatewayService

	tflog.Info(ctx, "Configured Trino-gateway client", map[string]any{"success": true})
}

func (p *TrinoGatewayProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewClusterResource,
	}
}

func (p *TrinoGatewayProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewClustersDataSource,
	}
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &TrinoGatewayProvider{
			version: version,
		}
	}
}
