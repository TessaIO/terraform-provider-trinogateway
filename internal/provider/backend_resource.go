// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/TessaIO/terraform-provider-trinogateway/internal/trinogateway"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &backendResource{}
	_ resource.ResourceWithConfigure   = &backendResource{}
	_ resource.ResourceWithImportState = &backendResource{}
)

// backendResourceModel maps Trino backend schema data.
type backendResourceModel struct {
	Name         types.String `tfsdk:"name"`
	ProxyTo      types.String `tfsdk:"proxy_to"`
	Active       types.Bool   `tfsdk:"active"`
	RoutingGroup types.String `tfsdk:"routing_group"`
	ExternalURL  types.String `tfsdk:"external_url"`
}

// NewBackendResource is a helper function to simplify the provider implementation.
func NewBackendResource() resource.Resource {
	return &backendResource{}
}

// backendResource is the resource implementation.
type backendResource struct {
	trinoGateway *trinogateway.TrinoGateway
}

// Metadata returns the resource type name.
func (r *backendResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backend"
}

// Schema defines the schema for the resource.
func (r *backendResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages a Trino backend.",
		MarkdownDescription: "Manages a Trino backend.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description:         "The name of the backend.",
				MarkdownDescription: "The name of the backend.",
				Required:            true,
			},
			"proxy_to": schema.StringAttribute{
				Description:         "The proxy URL for the backend.",
				MarkdownDescription: "The proxy URL for the backend.",
				Required:            true,
			},
			"active": schema.BoolAttribute{
				Description:         "Whether the backend is active.",
				MarkdownDescription: "Whether the backend is active.",
				Required:            true,
			},
			"routing_group": schema.StringAttribute{
				Description:         "The routing group for the backend.",
				MarkdownDescription: "The routing group for the backend.",
				Required:            true,
			},
			"external_url": schema.StringAttribute{
				Description:         "The external URL for the backend.",
				MarkdownDescription: "The external URL for the backend.",
				Computed:            true,
				Optional:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *backendResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.trinoGateway = client
	tflog.Debug(ctx, "Assigned TrinoGateway Client to the resource", map[string]any{"success": true})
}

// Create creates the resource and sets the initial Terraform state.
func (r *backendResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan backendResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	backend := trinogateway.CreateBackendRequest{
		Name:         plan.Name.ValueString(),
		ProxyTo:      plan.ProxyTo.ValueString(),
		Active:       plan.Active.ValueBool(),
		RoutingGroup: plan.RoutingGroup.ValueString(),
		ExternalURL:  plan.ExternalURL.ValueString(),
	}

	// Create new backend
	createdBackend, err := r.trinoGateway.CreateBackend(ctx, backend)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating backend",
			"Could not create backend, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "Backend created successfully in Trino gateway", map[string]any{"success": true, "createdBackend": createdBackend})

	// Map response body to schema and populate Computed attribute values
	plan.Name = types.StringValue(createdBackend.Name)
	plan.ProxyTo = types.StringValue(createdBackend.ProxyTo)
	plan.Active = types.BoolValue(createdBackend.Active)
	plan.RoutingGroup = types.StringValue(createdBackend.RoutingGroup)
	plan.ExternalURL = types.StringValue(createdBackend.ExternalURL)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "backend creation successfully saved in the plan state", map[string]any{"success": true})
}

// Read refreshes the Terraform state with the latest data.
func (r *backendResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state backendResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed backend value from trino gateway API
	backend, err := r.trinoGateway.GetBackend(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading TrinoGateway Backend",
			"Could not read TrinoGateway Backend Name "+state.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "GetBackend response from TrinoGateway API", map[string]any{"success": true, "backend": backend})

	// Map response body to schema and populate Computed attribute values
	state.Name = types.StringValue(backend.Name)
	state.ProxyTo = types.StringValue(backend.ProxyTo)
	state.Active = types.BoolValue(backend.Active)
	state.RoutingGroup = types.StringValue(backend.RoutingGroup)
	state.ExternalURL = types.StringValue(backend.ExternalURL)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "backend read successfully and save into the state from TrinoGateway API", map[string]any{"success": true})
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *backendResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan backendResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Backend will be updated", map[string]any{"backendName": plan.Name.String()})

	// Generate API request body from plan
	backend := trinogateway.UpdateBackendRequest{
		Name:         plan.Name.ValueString(),
		ProxyTo:      plan.ProxyTo.ValueStringPointer(),
		Active:       plan.Active.ValueBoolPointer(),
		RoutingGroup: plan.RoutingGroup.ValueStringPointer(),
		ExternalURL:  plan.ExternalURL.ValueStringPointer(),
	}

	// Update new backend
	err := r.trinoGateway.UpdateBackend(ctx, backend)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating backend",
			fmt.Sprintf("Could not update backend %q, unexpected error: %s", backend.Name, err),
		)
		return
	}

	tflog.Debug(ctx, "Backend updated successfully in Trino gateway", map[string]any{"success": true, "updatedBackend": backend})

	// Fetch updated Backend
	updatedBackend, err := r.trinoGateway.GetBackend(ctx, backend.Name)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Backend",
			"Could not read TrinoGateway Backend name "+plan.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.Name = types.StringValue(updatedBackend.Name)
	plan.ProxyTo = types.StringValue(updatedBackend.ProxyTo)
	plan.Active = types.BoolValue(updatedBackend.Active)
	plan.RoutingGroup = types.StringValue(updatedBackend.RoutingGroup)
	plan.ExternalURL = types.StringValue(updatedBackend.ExternalURL)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "backend read successfully and save into the state from TrinoGateway API", map[string]any{"success": true})
}

// Delete deletes the backend resource and removes the Terraform state on success.
func (r *backendResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from plan
	var state *backendResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Backend will be deleted %s", state.Name.ValueString()), map[string]any{"backendName": state.Active.ValueBool()})

	// Delete existing order
	if err := r.trinoGateway.DeleteBackend(ctx, state.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting TrinoGateway Backend",
			"Could not delete backend, unexpected error: "+err.Error(),
		)
		return

	}

	tflog.Debug(ctx, "backend deleted successfully from TrinoGateway API", map[string]any{"success": true})
}

// ImportState imports the state of an existing backend in the Trino gateway.
func (r *backendResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
