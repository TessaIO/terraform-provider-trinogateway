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
	resp.TypeName = req.ProviderTypeName + "_backend"
}

func (d *backendDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"clusters": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"proxy_to": schema.StringAttribute{
							Computed: true,
						},
						"active": schema.BoolAttribute{
							Computed: true,
						},
						"routing_group": schema.StringAttribute{
							Computed: true,
						},
						"external_url": schema.StringAttribute{
							Computed: true,
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
