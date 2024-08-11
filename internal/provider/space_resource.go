package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/qualitorque/terraform-provider-torque/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TorqueSpaceResource{}
var _ resource.ResourceWithImportState = &TorqueSpaceResource{}

func NewTorqueSpaceResource() resource.Resource {
	return &TorqueSpaceResource{}
}

// TorqueSpaceResource defines the resource implementation.
type TorqueSpaceResource struct {
	client *client.Client
}

// TorqueSpaceResourceModel describes the resource data model.
type TorqueSpaceResourceModel struct {
	Name  types.String `tfsdk:"space_name"`
	Color types.String `tfsdk:"color"`
	Icon  types.String `tfsdk:"icon"`
}

func (r *TorqueSpaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_space"
}

func (r *TorqueSpaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Creation of a new Torque space with associated entities (users, repos, etc...)",

		Attributes: map[string]schema.Attribute{
			"space_name": schema.StringAttribute{
				MarkdownDescription: "Space name to be create",
				Required:            true,
			},
			"color": schema.StringAttribute{
				MarkdownDescription: "Space color to be used for the new space",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("darkGreen"),
			},
			"icon": schema.StringAttribute{
				MarkdownDescription: "Space icon to be used",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("flow"),
			},
		},
	}
}

func (r *TorqueSpaceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *TorqueSpaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueSpaceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CreateSpace(data.Name.ValueString(), data.Color.ValueString(), data.Icon.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create space, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Resource Created Successful!")

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueSpaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueSpaceResourceModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	space, err := r.client.GetSpace(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading group details",
			"Could not read Torque group name "+data.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	// Treat HTTP 404 Not Found status as a signal to recreate resource
	// and return early
	if space.Name == "" {
		tflog.Error(ctx, "Space not found in Torque")
		resp.State.RemoveResource(ctx)
		return
	}

	data.Color = types.StringValue(space.Color)
	data.Icon = types.StringValue(space.Icon)

	// Set refreshed state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TorqueSpaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueSpaceResourceModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	err := r.client.UpdateSpace(data.Name.ValueString(), data.Color.ValueString(), data.Icon.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Torque space",
			"Could not update group, unexpected error: "+err.Error(),
		)
		return
	}

	space, err := r.client.GetSpace(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading space details",
			"Could not read Torque group name "+data.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	data.Color = types.StringValue(space.Color)
	data.Icon = types.StringValue(space.Icon)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TorqueSpaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueSpaceResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the space.
	err := r.client.DeleteSpace(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete space, got error: %s", err))
		return
	}

}

func (r *TorqueSpaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("space_name"), req, resp)
}
