// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/TessaIO/terraform-provider-trinogateway/internal/trinogateway"
)

var (
	_ datasource.DataSource              = &backendDataSource{}
	_ datasource.DataSourceWithConfigure = &backendDataSource{}
)

// backendDataSourceModel maps the data source schema data.
type backendDataSourceModel struct {
	Clusters []backendModel `tfsdk:"clusters"`
}

// backendModel maps Trino backend schema data.
type backendModel struct {
	Name         types.String `tfsdk:"name"`
	ProxyTo      types.String `tfsdk:"proxy_to"`
	Active       types.Bool   `tfsdk:"active"`
	RoutingGroup types.String `tfsdk:"routing_group"`
	ExternalURL  types.String `tfsdk:"external_url"`
}

func NewBackendDataSource() datasource.DataSource {
	return &backendDataSource{}
}

type backendDataSource struct {
	trinoGateway *trinogateway.TrinoGateway
}

func (d *backendDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backends"
}

func (d *backendDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "List of Trino backends.",
		MarkdownDescription: "List of Trino backends.",
		Attributes: map[string]schema.Attribute{
			"clusters": schema.ListNestedAttribute{
				Description:         "List of Trino backends.",
				MarkdownDescription: "List of Trino backends.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description:         "The name of the backend.",
							MarkdownDescription: "The name of the backend.",
							Computed:            true,
						},
						"proxy_to": schema.StringAttribute{
							Description:         "The proxy URL for the backend.",
							MarkdownDescription: "The proxy URL for the backend.",
							Computed:            true,
						},
						"active": schema.BoolAttribute{
							Description:         "Whether the backend is active.",
							MarkdownDescription: "Whether the backend is active.",
							Computed:            true,
						},
						"routing_group": schema.StringAttribute{
							Description:         "The routing group for the backend.",
							MarkdownDescription: "The routing group for the backend.",
							Computed:            true,
						},
						"external_url": schema.StringAttribute{
							Description:         "The external URL for the backend.",
							MarkdownDescription: "The external URL for the backend.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *backendDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*trinogateway.TrinoGateway)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *trinogateway.TrinoGateway, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.trinoGateway = client

	tflog.Debug(ctx, "Assigned TrinoGateway Client to the datasource", map[string]any{"success": true})
}

// Read reads the provided configuration and returns a response that contains the list of backends.
func (d *backendDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state backendDataSourceModel

	backends, err := d.trinoGateway.ListBackends(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read TrinoGateway Active Backends",
			err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "Successfully fetched all trino gateway active backends", map[string]any{"clusters": backends})

	// Map response body to model
	for _, backend := range backends {
		backendState := backendModel{
			Name:         types.StringValue(backend.Name),
			ProxyTo:      types.StringValue(backend.ProxyTo),
			Active:       types.BoolValue(backend.Active),
			RoutingGroup: types.StringValue(backend.RoutingGroup),
			ExternalURL:  types.StringValue(backend.ExternalURL),
		}

		state.Clusters = append(state.Clusters, backendState)
	}

	tflog.Debug(ctx, "Successfully constructed the clusters state", map[string]any{"clusters": state.Clusters})
	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
