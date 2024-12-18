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
var _ resource.Resource = &TorqueWorkflowResource{}
var _ resource.ResourceWithImportState = &TorqueWorkflowResource{}

func NewTorqueWorkflowResource() resource.Resource {
	return &TorqueWorkflowResource{}
}

// TorqueWorkflowResource defines the resource implementation.
type TorqueWorkflowResource struct {
	client *client.Client
}

// TorqueWorkflowResourceModel describes the resource data model.
type TorqueWorkflowResourceModel struct {
	Name          types.String `tfsdk:"name"`
	SpaceName     types.String `tfsdk:"space_name"`
	RepoName      types.String `tfsdk:"repository_name"`
	LaunchAllowed types.String `tfsdk:"launch_allowed"`
}

func (r *TorqueWorkflowResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_workflow"
}

func (r *TorqueWorkflowResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Enables an existing Torque workflow so it will be allowed to launch and executed.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the workflow to enable.",
				Optional:            false,
				Computed:            false,
				Required:            true,
			},
			"space_name": schema.StringAttribute{
				MarkdownDescription: "Space the workflow belongs to",
				Optional:            false,
				Computed:            false,
				Required:            true,
			},
			"repository_name": schema.StringAttribute{
				MarkdownDescription: "Repository where the workflow source code is",
				Optional:            false,
				Computed:            false,
				Required:            true,
			},
			"launch_allowed": schema.StringAttribute{
				MarkdownDescription: "Indicates whether this workflow is enabled and allowed to be launched",
				Optional:            false,
				Computed:            false,
				Required:            true,
			},
		},
	}
}

func (r *TorqueWorkflowResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueWorkflowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueWorkflowResourceModel

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

func (r *TorqueWorkflowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueWorkflowResourceModel

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

func (r *TorqueWorkflowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state TorqueWorkflowResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	diags := req.Plan.Get(ctx, &plan)

	current_space := state.Name

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing order
	err := r.client.UpdateSpace(current_space.ValueString(), plan.Name.ValueString(), plan.Color.ValueString(), plan.Icon.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Torque space",
			"Could not update Torque Space name, unexpected error: "+err.Error(),
		)
		return
	}

	space, err := r.client.GetSpace(plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading space details",
			"Could not read Torque Space name "+plan.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	plan.Color = types.StringValue(space.Color)
	plan.Icon = types.StringValue(space.Icon)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TorqueWorkflowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueWorkflowResourceModel

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

func (r *TorqueWorkflowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("space_name"), req, resp)
}
