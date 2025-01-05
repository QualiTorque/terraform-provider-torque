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
var _ resource.Resource = &TorqueSpaceWorkflowResource{}
var _ resource.ResourceWithImportState = &TorqueSpaceWorkflowResource{}

func NewTorqueSpaceWorkflowResource() resource.Resource {
	return &TorqueSpaceWorkflowResource{}
}

// TorqueSpaceWorkflowResource defines the resource implementation.
type TorqueSpaceWorkflowResource struct {
	client *client.Client
}

// TorqueSpaceWorkflowResourceModel describes the resource data model.
type TorqueSpaceWorkflowResourceModel struct {
	Name          types.String `tfsdk:"name"`
	SpaceName     types.String `tfsdk:"space_name"`
	RepoName      types.String `tfsdk:"repository_name"`
	LaunchAllowed types.Bool   `tfsdk:"launch_allowed"`
	CustomIcon    types.String `tfsdk:"custom_icon"`
}

func (r *TorqueSpaceWorkflowResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_space_workflow"
}

func (r *TorqueSpaceWorkflowResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Allows to existing Torque workflow with a space scope so it will be allowed to published and displayed in the self-service catalog.",

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
			"custom_icon": schema.StringAttribute{
				MarkdownDescription: "Custom icon key to associate with this catalog item. The key can be referenced from a torque_space_custom_icon key attribute.",
				Required:            false,
				Optional:            true,
				Computed:            false,
			},
		},
	}
}

func (r *TorqueSpaceWorkflowResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueSpaceWorkflowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueSpaceWorkflowResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.PublishBlueprintInSpace(data.SpaceName.ValueString(), data.RepoName.ValueString(), data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to publish workflow to self-service catalog, got error: %s", err))
		return
	}
	if !data.CustomIcon.IsNull() {
		err := r.client.SetCatalogItemCustomIcon(data.SpaceName.ValueString(), data.Name.ValueString(), data.RepoName.ValueString(), data.CustomIcon.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Catalog Item, failed to set catalog item custom icon, got error: %s", err))
			return
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueSpaceWorkflowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

}

func (r *TorqueSpaceWorkflowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueSpaceWorkflowResourceModel
	var state TorqueSpaceWorkflowResourceModel
	const default_icon = "nodes"

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if data.CustomIcon.IsNull() {
		err := r.client.SetCatalogItemIcon(data.SpaceName.ValueString(), data.Name.ValueString(), data.RepoName.ValueString(), default_icon)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update workflow custom icon, failed to set catalog item custom icon, got error: %s", err))
			return
		}
	} else {
		if data.CustomIcon.ValueString() != state.CustomIcon.ValueString() {
			err := r.client.SetCatalogItemCustomIcon(data.SpaceName.ValueString(), data.Name.ValueString(), data.RepoName.ValueString(), data.CustomIcon.ValueString())
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update workflow custom icon, failed to set catalog item custom icon, got error: %s", err))
				return
			}
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueSpaceWorkflowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueSpaceWorkflowResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UnpublishBlueprintInSpace(data.SpaceName.ValueString(), data.RepoName.ValueString(), data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to unpublish workflow to self-service catalog, got error: %s", err))
		return
	}

}

func (r *TorqueSpaceWorkflowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
