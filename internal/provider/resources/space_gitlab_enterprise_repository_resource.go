package resources

import (
	"context"
	"fmt"
	"strings"

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
var _ resource.Resource = &TorqueSpaceGitlabEnterpriseRepositoryResource{}
var _ resource.ResourceWithImportState = &TorqueSpaceGitlabEnterpriseRepositoryResource{}

func NewTorqueSpaceGitlabEnterpriseRepositoryResource() resource.Resource {
	return &TorqueSpaceGitlabEnterpriseRepositoryResource{}
}

// TorqueAgentSpaceAssociationResource defines the resource implementation.
type TorqueSpaceGitlabEnterpriseRepositoryResource struct {
	client *client.Client
}

type TorqueSpaceGitlabEnterpriseRepositoryResourceModel struct {
	SpaceName      types.String `tfsdk:"space_name"`
	RepositoryName types.String `tfsdk:"repository_name"`
	RepositoryUrl  types.String `tfsdk:"repository_url"`
	Token          types.String `tfsdk:"token"`
	Branch         types.String `tfsdk:"branch"`
	CredentialName types.String `tfsdk:"credential_name"`
	UseAllAgents   types.Bool   `tfsdk:"use_all_agents"`
	Agents         types.List   `tfsdk:"agents"`
}

func (r *TorqueSpaceGitlabEnterpriseRepositoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_gitlab_enterprise_repository_space_association"
}

func (r *TorqueSpaceGitlabEnterpriseRepositoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Onboard a new GitlabEnterprise repository into an existing space",

		Attributes: map[string]schema.Attribute{
			"space_name": schema.StringAttribute{
				MarkdownDescription: "Existing Torque Space name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"repository_name": schema.StringAttribute{
				Description: "The name of the GitlabEnterprise repository to onboard. In this example, repo_name",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"repository_url": schema.StringAttribute{
				Description: "The url of the specific GitlabEnterprise repository/project to onboard. For example: https://gitlab-on-prem.example.com/repo_name",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"token": schema.StringAttribute{
				Description: "Authentication Token to the project/repository. If omitted, existing credentials provided in the credential_name field will be used for authentication. If provided, a new credentials object will be created.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				DeprecationMessage: "The token attribute is deprecated and will be removed in a future release. Use the torque_git_credentials resource to store the token and reference it in this resource using the credential_name attribute instead.",
			},
			"branch": schema.StringAttribute{
				Description: "Repository branch to use for blueprints and automation assets",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"credential_name": schema.StringAttribute{
				Description: "The name of the Credentials to use/create. Must be unique in the space.",
				Required:    true,
			},
			"use_all_agents": schema.BoolAttribute{
				Description: "Whether all associated agents can be used to onboard and sync this repository. Must be set to false if agents attribute is used.",
				Default:     booldefault.StaticBool(true),
				Optional:    true,
				Computed:    true,
				Validators:  []validator.Bool{UseAllAgentsValidator{}},
			},
			"agents": schema.ListAttribute{
				Description: "List of specific agents to use to onboard and sync this repository. Cannot be specified when use_all_agents is true.",
				Required:    false,
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *TorqueSpaceGitlabEnterpriseRepositoryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueSpaceGitlabEnterpriseRepositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueSpaceGitlabEnterpriseRepositoryResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	agents := []string{}
	if !data.Agents.IsNull() {
		for _, agent := range data.Agents.Elements() {
			agents = append(agents, strings.Trim(agent.String(), "\""))
		}
	}
	err := r.client.OnboardGitlabEnterpriseRepoToSpace(data.SpaceName.ValueString(), data.RepositoryName.ValueString(),
		data.RepositoryUrl.ValueString(), data.Token.ValueStringPointer(), data.Branch.ValueString(), data.CredentialName.ValueString(), agents, data.UseAllAgents.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to onboard repository to space, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Resource Created Successful!")

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueSpaceGitlabEnterpriseRepositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueSpaceGitlabEnterpriseRepositoryResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.

	// Save updated data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueSpaceGitlabEnterpriseRepositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueSpaceGitlabEnterpriseRepositoryResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	agents := []string{}
	if !data.Agents.IsNull() {
		for _, agent := range data.Agents.Elements() {
			agents = append(agents, strings.Trim(agent.String(), "\""))
		}
	}
	err := r.client.UpdateRepoConfiguration(data.SpaceName.ValueString(), data.RepositoryName.ValueString(),
		data.CredentialName.ValueString(), agents, data.UseAllAgents.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update repository configuration, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueSpaceGitlabEnterpriseRepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueSpaceGitlabEnterpriseRepositoryResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Remove repo from space.
	err := r.client.RemoveRepoFromSpace(data.SpaceName.ValueString(), data.RepositoryName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to remove repository from space, got error: %s", err))
		return
	}

}

func (r *TorqueSpaceGitlabEnterpriseRepositoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

type UseAllAgentsValidator struct{}

func (v UseAllAgentsValidator) Description(ctx context.Context) string {
	return "Ensures use_all_agents is false when agents are provided."
}

func (v UseAllAgentsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v UseAllAgentsValidator) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	useAllAgents := req.ConfigValue.ValueBool()
	var agents []types.String

	// Fetch the agents attribute
	if diags := req.Config.GetAttribute(ctx, path.Root("agents"), &agents); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Check if use_all_agents is true and agents should be empty
	if useAllAgents && len(agents) > 0 {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"Cannot specify agents when use_all_agents is true.",
		)
		return
	}

	// Check if use_all_agents is false and agents list must have at least 1 element
	if !useAllAgents && len(agents) == 0 {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"Agents list must contain at least one element when use_all_agents is false.",
		)
	}
}
