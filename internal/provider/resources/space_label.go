package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/qualitorque/terraform-provider-torque/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TorqueSpaceLabelResource{}
var _ resource.ResourceWithImportState = &TorqueSpaceLabelResource{}

func NewTorqueSpaceLabelResource() resource.Resource {
	return &TorqueSpaceLabelResource{}
}

// TorqueSpaceLabelResource defines the resource implementation.
type TorqueSpaceLabelResource struct {
	client *client.Client
}

// TorqueSpaceLabelResourceModel describes the resource data model.
type TorqueSpaceLabelResourceModel struct {
	SpaceName   types.String `tfsdk:"space_name"`
	Name        types.String `tfsdk:"name"`
	Color       types.String `tfsdk:"color"`
	QuickFilter types.Bool   `tfsdk:"quick_filter"`
}

func (r *TorqueSpaceLabelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_space_label"
}

func (r *TorqueSpaceLabelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Creation of a new label to be used in Torque space and can be associated to catalog items.",

		Attributes: map[string]schema.Attribute{
			"space_name": schema.StringAttribute{
				MarkdownDescription: "Name of the space where this label will be added to Torque",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the new label to be added to torque",
				Required:            true,
			},
			"color": schema.StringAttribute{
				MarkdownDescription: "Color of the label. Allowed values: aws, darkGray, frogGreen, pink, orange, blueGray, blue, bordeaux, teal, grey",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"aws", "darkGray", "frogGreen", "pink", "orange", "blueGray", "blue", "bordeaux", "teal", "grey"}...),
				},
			},
			"quick_filter": schema.BoolAttribute{
				MarkdownDescription: "Display this label as a quick filter in the self-service catalog.",
				Required:            false,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
		},
	}
}

func (r *TorqueSpaceLabelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueSpaceLabelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueSpaceLabelResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CreateLabel(data.SpaceName.ValueString(),
		data.Name.ValueString(), data.Color.ValueString(), data.QuickFilter.ValueBool())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create label, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Resource Created Successful!")

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueSpaceLabelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueSpaceLabelResourceModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	label, err := r.client.GetLabel(data.SpaceName.ValueString(), data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading label details",
			"Could not read Torque label name "+data.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	// Treat HTTP 404 Not Found status as a signal to recreate resource
	// and return early
	if label.Name == "" {
		tflog.Error(ctx, "label not found in Torque")
		resp.State.RemoveResource(ctx)
		return
	}

	data.Name = types.StringValue(label.Name)
	data.Color = types.StringValue(label.Color)

	// Set refreshed state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TorqueSpaceLabelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state TorqueSpaceLabelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	current_name := state.Name
	// Update existing order
	err := r.client.UpdateLabel(current_name.ValueString(), data.SpaceName.ValueString(), data.Name.ValueString(), data.Color.ValueString(), data.QuickFilter.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Torque label",
			"Could not update group, unexpected error: "+err.Error(),
		)
		return
	}

	label, err := r.client.GetLabel(data.SpaceName.ValueString(), data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading label details",
			"Could not read Torque label name "+data.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	data.Name = types.StringValue(label.Name)
	data.Color = types.StringValue(label.Color)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TorqueSpaceLabelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueSpaceLabelResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the space.
	err := r.client.DeleteLabel(data.SpaceName.ValueString(), data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete label, got error: %s", err))
		return
	}

}

func (r *TorqueSpaceLabelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
