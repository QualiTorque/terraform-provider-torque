package provider

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
var _ resource.Resource = &TorqueSpaceCodeCommitRepositoryResource{}
var _ resource.ResourceWithImportState = &TorqueSpaceCodeCommitRepositoryResource{}

func NewTorqueSpaceCodeCommitRepositoryResource() resource.Resource {
	return &TorqueSpaceCodeCommitRepositoryResource{}
}

// TorqueAgentSpaceAssociationResource defines the resource implementation.
type TorqueSpaceCodeCommitRepositoryResource struct {
	client *client.Client
}

type TorqueSpaceCodeCommitRepositoryResourceModel struct {
	SpaceName      types.String `tfsdk:"space_name"`
	RepositoryUrl  types.String `tfsdk:"repository_url"`
	RoleArn        types.String `tfsdk:"role_arn"`
	AwsRegion      types.String `tfsdk:"aws_region"`
	ExternalId     types.String `tfsdk:"external_id"`
	Username       types.String `tfsdk:"username"`
	Password       types.String `tfsdk:"password"`
	Branch         types.String `tfsdk:"branch"`
	RepositoryName types.String `tfsdk:"repository_name"`
}

func (r *TorqueSpaceCodeCommitRepositoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_codecommit_repository_space_association"
}

func (r *TorqueSpaceCodeCommitRepositoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Onboard a new repository into an existing space",

		Attributes: map[string]schema.Attribute{
			"space_name": schema.StringAttribute{
				MarkdownDescription: "Existing Torque Space name",
				Required:            true,
			},
			"repository_url": schema.StringAttribute{
				Description: "Repository URL. For example: https://github.com/<org>/<repo>",
				Required:    true,
			},
			"role_arn": schema.StringAttribute{
				Description: "Repository type. Available types: github, bitbucket, gitlab, azure (for Azure DevOps)",
				Required:    true,
			},
			"aws_region": schema.StringAttribute{
				Description: "Repository type. Available types: github, bitbucket, gitlab, azure (for Azure DevOps)",
				Required:    true,
			},
			"external_id": schema.StringAttribute{
				Description: "Repository type. Available types: github, bitbucket, gitlab, azure (for Azure DevOps)",
				Required:    true,
			},
			"username": schema.StringAttribute{
				Description: "Repository type. Available types: github, bitbucket, gitlab, azure (for Azure DevOps)",
				Required:    true,
			},
			"password": schema.StringAttribute{
				Description: "Repository type. Available types: github, bitbucket, gitlab, azure (for Azure DevOps)",
				Required:    true,
			},
			"branch": schema.StringAttribute{
				Description: "Repository branch to use for blueprints and automation assets",
				Optional:    true,
			},
			"repository_name": schema.StringAttribute{
				Description: "The name of the repository to onboard in the newly created space",
				Required:    true,
			},
		},
	}
}

func (r *TorqueSpaceCodeCommitRepositoryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueSpaceCodeCommitRepositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueSpaceCodeCommitRepositoryResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.OnboardCodeCommitRepoToSpace(data.SpaceName.ValueString(), data.RepositoryName.ValueString(), data.RoleArn.ValueString(),
		data.RepositoryUrl.ValueString(), data.AwsRegion.ValueString(), data.Branch.ValueString(), data.ExternalId.ValueString(), data.Username.ValueString(), data.Password.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to onboard repository to space, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Resource Created Successful!")

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueSpaceCodeCommitRepositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueSpaceCodeCommitRepositoryResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueSpaceCodeCommitRepositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueSpaceCodeCommitRepositoryResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueSpaceCodeCommitRepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueSpaceCodeCommitRepositoryResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Remove repo from space.
	err := r.client.RemoveRepoFromSpace(data.SpaceName.ValueString(), data.RepositoryName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to attach agent to space, got error: %s", err))
		return
	}

}

func (r *TorqueSpaceCodeCommitRepositoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
