package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/qualitorque/terraform-provider-torque/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TorqueSpaceGitCredentialsResource{}
var _ resource.ResourceWithImportState = &TorqueSpaceGitCredentialsResource{}

func NewTorqueSpaceGitCredentialsResource() resource.Resource {
	return &TorqueSpaceGitCredentialsResource{}
}

// TorqueSpaceGitCredentialsResource defines the resource implementation.
type TorqueSpaceGitCredentialsResource struct {
	client *client.Client
}

// TorqueSpaceGitCredentialsResourceModel describes the resource data model.
type TorqueSpaceGitCredentialsResourceModel struct {
	SpaceName   types.String `tfsdk:"space_name"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Token       types.String `tfsdk:"token"`
	Type        types.String `tfsdk:"type"`
	CloudType   types.String `tfsdk:"cloudtype"`
}

func (r *TorqueSpaceGitCredentialsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_space_git_credentials"
}

func (r *TorqueSpaceGitCredentialsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Creation of a new git credentials resource in a specific space, which can later be used to onboard a git repository. Supported repositories are github, gitlab enterprise, azure devops and bitbucket.",

		Attributes: map[string]schema.Attribute{
			"space_name": schema.StringAttribute{
				MarkdownDescription: "Name of the space to create the credentials in.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the credentials.",
				Required:            true,
				Computed:            false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the credentials.",
				Required:            true,
				Computed:            false,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "Access token the credentials will use.",
				Required:            true,
				Computed:            false,
				Sensitive:           true,
			},
			"cloudtype": schema.StringAttribute{
				MarkdownDescription: "Credentials type identifier",
				Required:            false,
				Computed:            true,
				Default:             stringdefault.StaticString("sourceControl"),
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of git repository these credentials are for. Supported types are github, bitbucket, azureDevops and gitlabEnterprise.",
				Required:            true,
				Computed:            false,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"github", "bitbucket", "azureDevops", "gitlabEnterprise"}...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *TorqueSpaceGitCredentialsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueSpaceGitCredentialsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueSpaceGitCredentialsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.CreateSpaceCredentials(data.SpaceName.ValueString(), data.Name.ValueString(), data.Description.ValueString(), data.CloudType.ValueString(), data.Type.ValueString(), data.Token.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create space git credentials, got error: %s", err))
		return
	}
	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *TorqueSpaceGitCredentialsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueSpaceGitCredentialsResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

}

func (r *TorqueSpaceGitCredentialsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueSpaceGitCredentialsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.UpdateSpaceCredentials(data.SpaceName.ValueString(), data.Name.ValueString(), data.Description.ValueString(), data.CloudType.ValueString(), data.Type.ValueString(), data.Token.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update space git credentials, got error: %s", err))
		return
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueSpaceGitCredentialsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueSpaceGitCredentialsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.DeleteSpaceCredentials(data.SpaceName.ValueString(), data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete space git credentials, got error: %s", err))
		return
	}
}

func (r *TorqueSpaceGitCredentialsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
