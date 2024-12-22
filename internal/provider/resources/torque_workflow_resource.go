package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	LaunchAllowed types.Bool   `tfsdk:"launch_allowed"`
	SelfService   types.Bool   `tfsdk:"self_service"`
}

func (r *TorqueWorkflowResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_workflow"
}

func (r *TorqueWorkflowResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Allows to enable and publish existing Torque workflow so it will be allowed to be executed and displayed in the self-service catalog.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the workflow to enable.",
				Optional:            false,
				Computed:            false,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				}},
			"space_name": schema.StringAttribute{
				MarkdownDescription: "Space the workflow belongs to",
				Optional:            false,
				Computed:            false,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"repository_name": schema.StringAttribute{
				MarkdownDescription: "Repository where the workflow source code is",
				Optional:            false,
				Computed:            false,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"launch_allowed": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether this workflow is enabled and allowed to be launched",
				Optional:            false,
				Computed:            true,
				Required:            false,
				Default:             booldefault.StaticBool(true),
			},
			"self_service": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether this workflow is displayed in the self-service catalog. For workflows with Space scope, then this field can be omitted and will always be true.",
				Optional:            true,
				Required:            false,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
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

	if data.SelfService.ValueBool() {
		err := r.client.PublishBlueprintInSpace(data.SpaceName.ValueString(), data.RepoName.ValueString(), data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to publish workflow to self-service catalog, got error: %s", err))
			return
		}
	}
	err := r.client.AllowLaunch(data.Name.ValueString(), data.RepoName.ValueString(), data.SpaceName.ValueString(), data.LaunchAllowed.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to publish workflow to self-service catalog, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueWorkflowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

}

func (r *TorqueWorkflowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueWorkflowResourceModel
	var state TorqueWorkflowResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	if data.SelfService.ValueBool() {
		err := r.client.PublishBlueprintInSpace(data.SpaceName.ValueString(), data.RepoName.ValueString(), data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to publish workflow to self-service catalog, got error: %s", err))
			return
		}
	} else {
		err := r.client.UnpublishBlueprintInSpace(data.SpaceName.ValueString(), data.RepoName.ValueString(), data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to unpublish workflow to self-service catalog, got error: %s", err))
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueWorkflowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueWorkflowResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.LaunchAllowed = types.BoolValue(false)
	err := r.client.AllowLaunch(data.Name.ValueString(), data.RepoName.ValueString(), data.SpaceName.ValueString(), data.LaunchAllowed.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to disable workflow, got error: %s", err))
		return
	}
	err = r.client.UnpublishBlueprintInSpace(data.SpaceName.ValueString(), data.RepoName.ValueString(), data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to unpublish workflow to self-service catalog, got error: %s", err))
		return
	}

}

func (r *TorqueWorkflowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
