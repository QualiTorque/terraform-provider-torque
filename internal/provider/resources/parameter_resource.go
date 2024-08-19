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
var _ resource.Resource = &TorqueParameterResource{}
var _ resource.ResourceWithImportState = &TorqueParameterResource{}

func NewTorqueParameterResource() resource.Resource {
	return &TorqueParameterResource{}
}

// TorqueParameterResource defines the resource implementation.
type TorqueParameterResource struct {
	client *client.Client
}

// TorqueParameterResourceModel describes the resource data model.
type TorqueParameterResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Value       types.String `tfsdk:"value"`
	Sensitive   types.Bool   `tfsdk:"sensitive"`
	Description types.String `tfsdk:"description"`
}

func (r *TorqueParameterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_parameter"
}

func (r *TorqueParameterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Creation of a new parameter is a Torque",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the new parameter to be added to torque",
				Required:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "Tag value to be set as the parameter",
				Optional:            true,
				Computed:            false,
			},
			"sensitive": schema.BoolAttribute{
				MarkdownDescription: "Sensitive or not",
				Optional:            true,
				Computed:            false,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Parameter description",
				Optional:            true,
				Computed:            false,
			},
		},
	}
}

func (r *TorqueParameterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueParameterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueParameterResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.AddAccountParameter(data.Name.ValueString(),
		data.Value.ValueString(), data.Sensitive.ValueBool(), data.Description.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create parameter, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Resource Created Successful!")

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueParameterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueParameterResourceModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	parameter, err := r.client.GetAccountParameter(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading group details",
			"Could not read Torque group name "+data.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	// Treat HTTP 404 Not Found status as a signal to recreate resource
	// and return early
	if parameter.Name == "" {
		tflog.Error(ctx, "Parameter not found in Torque")
		resp.State.RemoveResource(ctx)
		return
	}

	data.Description = types.StringValue(parameter.Description)
	data.Value = types.StringValue(parameter.Value)
	data.Sensitive = types.BoolValue(parameter.Sensitive)

	// Set refreshed state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TorqueParameterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueParameterResourceModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	err := r.client.UpdateAccountParameter(data.Name.ValueString(), data.Value.ValueString(), data.Sensitive.ValueBool(), data.Description.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Torque parameter",
			"Could not update group, unexpected error: "+err.Error(),
		)
		return
	}

	param, err := r.client.GetAccountParameter(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading parameter details",
			"Could not read Torque parameter name "+data.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	data.Description = types.StringValue(param.Description)
	data.Sensitive = types.BoolValue(param.Sensitive)
	data.Value = types.StringValue(param.Value)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TorqueParameterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueParameterResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the space.
	err := r.client.DeleteAccountParameter(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete parameter, got error: %s", err))
		return
	}

}

func (r *TorqueParameterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
