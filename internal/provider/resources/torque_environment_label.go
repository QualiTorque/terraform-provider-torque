package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/qualitorque/terraform-provider-torque/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TorqueEnvironmentLabelResource{}
var _ resource.ResourceWithImportState = &TorqueEnvironmentLabelResource{}

func NewTorqueEnvironmentLabelResource() resource.Resource {
	return &TorqueEnvironmentLabelResource{}
}

// TorqueEnvironmentLabelResource defines the resource implementation.
type TorqueEnvironmentLabelResource struct {
	client *client.Client
}

// TorqueEnvironmentLabelResourceModel describes the resource data model.
type TorqueEnvironmentLabelResourceModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

func (r *TorqueEnvironmentLabelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_environment_label"
}

func (r *TorqueEnvironmentLabelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Creation of a new label to be used in Torque environment and can be associated to catalog items.",

		Attributes: map[string]schema.Attribute{
			"key": schema.StringAttribute{
				MarkdownDescription: "Value of the new environment label to be added to torque",
				Required:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "Value of the label",
				Required:            true,
			},
		},
	}
}

func (r *TorqueEnvironmentLabelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueEnvironmentLabelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueEnvironmentLabelResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CreateEnvironmentLabel(data.Key.ValueString(), data.Value.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create environment label, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Resource Created Successful!")

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueEnvironmentLabelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueEnvironmentLabelResourceModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	label, err := r.client.GetEnvironmentLabel(data.Key.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading label details",
			"Could not read Torque label name "+data.Value.ValueString()+": "+err.Error(),
		)
		return
	}

	// Treat HTTP 404 Not Found status as a signal to recreate resource
	// and return early
	if label.Value == "" {
		tflog.Error(ctx, "label not found in Torque")
		resp.State.RemoveResource(ctx)
		return
	}
	data.Key = types.StringValue(label.Key)
	data.Value = types.StringValue(label.Value)
	// Set refreshed state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TorqueEnvironmentLabelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state TorqueEnvironmentLabelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateEnvironmentLabel(state.Key.ValueString(), state.Value.ValueString(), data.Key.ValueString(), data.Value.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Torque label",
			"Could not update label, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TorqueEnvironmentLabelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueEnvironmentLabelResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the environment label.
	err := r.client.DeleteEnvironmentLabel(data.Key.ValueString(), data.Value.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete environment label, got error: %s", err))
		return
	}

}

func (r *TorqueEnvironmentLabelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
