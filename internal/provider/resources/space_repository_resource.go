package resources

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/qualitorque/terraform-provider-torque/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TorqueSpaceRepositoryResource{}
var _ resource.ResourceWithImportState = &TorqueSpaceRepositoryResource{}

func NewTorqueSpaceRepositoryResource() resource.Resource {
	return &TorqueSpaceRepositoryResource{}
}

// TorqueAgentSpaceAssociationResource defines the resource implementation.
type TorqueSpaceRepositoryResource struct {
	client *client.Client
}

type TorqueSpaceRepositoryResourceModel struct {
	SpaceName      types.String `tfsdk:"space_name"`
	RepoUrl        types.String `tfsdk:"repository_url"`
	RepoToken      types.String `tfsdk:"access_token"`
	RepoType       types.String `tfsdk:"repository_type"`
	RepoBranch     types.String `tfsdk:"branch"`
	RepoName       types.String `tfsdk:"repository_name"`
	CredentialName types.String `tfsdk:"credential_name"`
}

func (r *TorqueSpaceRepositoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_repository_space_association"
}

func (r *TorqueSpaceRepositoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Onboard a new repository into an existing space",

		Attributes: map[string]schema.Attribute{
			"space_name": schema.StringAttribute{
				MarkdownDescription: "Existing Torque Space name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"repository_url": schema.StringAttribute{
				Description: "Repository URL. For example: https://github.com/<org>/<repo>",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"access_token": schema.StringAttribute{
				Description: "Personal Access Token (PAT) to authenticate with to the repository. Credentials will be automatically created with the specified token, or use existing credentials instead.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				DeprecationMessage: "The access_token attribute is deprecated and will be removed in a future release. Use the torque_git_credentials resource to store the token and reference it in this resource using the credential_name attribute instead.",
				Validators: []validator.String{
					// Validate only this attribute or other_attr is configured or neither.
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("credential_name"),
					}...),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^\S.*$`),
						"must not be empty string",
					),
				},
			},
			"repository_type": schema.StringAttribute{
				Description: "Repository type. Available types: github, bitbucket, gitlab, azure (for Azure DevOps). For CodeCommit, Please use torque_codecommit_repository_space_association resource. For Gitlab Enterprise please use torque_gitlab_enterprise_repository_space_association resource",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"github", "bitbucket", "gitlab", "azure"}...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"branch": schema.StringAttribute{
				Description: "Repository branch to use for blueprints and automation assets",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"repository_name": schema.StringAttribute{
				Description: "The name of the repository to onboard in the newly created space",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"credential_name": schema.StringAttribute{
				Description: "The name of existing credentials to use.",
				Required:    false,
				Optional:    true,
				Validators: []validator.String{
					// Validate only this attribute or other_attr is configured or neither.
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("access_token"),
					}...),
				},
			},
		},
	}
}

func (r *TorqueSpaceRepositoryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueSpaceRepositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueSpaceRepositoryResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.OnboardRepoToSpace(data.SpaceName.ValueString(), data.RepoName.ValueString(), data.RepoType.ValueString(),
		data.RepoUrl.ValueString(), data.RepoToken.ValueStringPointer(), data.RepoBranch.ValueString(), data.CredentialName.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to onboard repository to space, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Resource Created Successful!")

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueSpaceRepositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueSpaceRepositoryResourceModel

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

func (r *TorqueSpaceRepositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueSpaceRepositoryResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateRepoCredentials(data.SpaceName.ValueString(), data.RepoName.ValueString(), data.CredentialName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to change repository credentials, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueSpaceRepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueSpaceRepositoryResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Remove repo from space.
	err := r.client.RemoveRepoFromSpace(data.SpaceName.ValueString(), data.RepoName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to attach agent to space, got error: %s", err))
		return
	}

}

func (r *TorqueSpaceRepositoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
