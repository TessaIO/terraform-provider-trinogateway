package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/TessaIO/terraform-provider-trinogateway/internal/trinogateway"
)

var (
	_ datasource.DataSource              = &clusterDataSource{}
	_ datasource.DataSourceWithConfigure = &clusterDataSource{}
)

// clusterDataSourceModel maps the data source schema data.
type clusterDataSourceModel struct {
	Clusters []clusterModel `tfsdk:"clusters"`
}

// clusterModel maps Trino cluster schema data.
type clusterModel struct {
	Name         types.String `tfsdk:"name"`
	ProxyTo      types.String `tfsdk:"proxy_to"`
	Active       types.Bool   `tfsdk:"active"`
	RoutingGroup types.String `tfsdk:"routing_group"`
	ExternalURL  types.String `tfsdk:"external_url"`
}

func NewClustersDataSource() datasource.DataSource {
	return &clusterDataSource{}
}

type clusterDataSource struct {
	trinoGateway *trinogateway.TrinoGateway
}

func (d *clusterDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

func (d *clusterDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
func (d *clusterDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
}

func (d *clusterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state clusterDataSourceModel

	clusters, err := d.trinoGateway.ListBackends(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read TrinoGateway Active Clusters",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, cluster := range clusters {
		clusterState := clusterModel{
			Name:         types.StringValue(cluster.Name),
			ProxyTo:      types.StringValue(cluster.ProxyTo),
			Active:       types.BoolValue(cluster.Active),
			RoutingGroup: types.StringValue(cluster.RoutingGroup),
			ExternalURL:  types.StringValue(cluster.ExternalURL),
		}

		state.Clusters = append(state.Clusters, clusterState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
