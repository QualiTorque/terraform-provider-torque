package resources

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/qualitorque/terraform-provider-torque/client"
	"github.com/qualitorque/terraform-provider-torque/internal/validators"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TorqueSpaceAdoServerRepositoryResource{}
var _ resource.ResourceWithImportState = &TorqueSpaceAdoServerRepositoryResource{}

func NewTorqueSpaceAdoServerRepositoryResource() resource.Resource {
	return &TorqueSpaceAdoServerRepositoryResource{}
}

// TorqueAgentSpaceAssociationResource defines the resource implementation.
type TorqueSpaceAdoServerRepositoryResource struct {
	client *client.Client
}

type TorqueSpaceAdoServerRepositoryResourceModel struct {
	SpaceName       types.String `tfsdk:"space_name"`
	RepositoryName  types.String `tfsdk:"repository_name"`
	RepositoryUrl   types.String `tfsdk:"repository_url"`
	Token           types.String `tfsdk:"token"`
	Branch          types.String `tfsdk:"branch"`
	CredentialName  types.String `tfsdk:"credential_name"`
	UseAllAgents    types.Bool   `tfsdk:"use_all_agents"`
	Agents          types.List   `tfsdk:"agents"`
	TimeOut         types.Int32  `tfsdk:"timeout"`
	AutoRegisterEac types.Bool   `tfsdk:"auto_register_eac"`
}

func (r *TorqueSpaceAdoServerRepositoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_ado_server_repository_space_association"
}

func (r *TorqueSpaceAdoServerRepositoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Onboard a new Ado Server repository into an existing space",

		Attributes: map[string]schema.Attribute{
			"space_name": schema.StringAttribute{
				MarkdownDescription: "Existing Torque Space name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"repository_name": schema.StringAttribute{
				Description: "The name of the Ado Server repository to onboard. In this example, repo_name",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"repository_url": schema.StringAttribute{
				Description: "The url of the specific Ado Server repository/project to onboard. For example: https://ado-on-prem.example.com/repo_name",
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
				Validators:  []validator.Bool{validators.UseAllAgentsValidator{}},
			},
			"agents": schema.ListAttribute{
				Description: "List of specific agents to use to onboard and sync this repository. Cannot be specified when use_all_agents is true.",
				Required:    false,
				Optional:    true,
				ElementType: types.StringType,
			},
			"timeout": schema.Int32Attribute{
				Description: "Time in minutes to wait for Torque to sync the repository during the onboarding. Default is 1 minute.",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Default:     int32default.StaticInt32(1),
			},
			"auto_register_eac": schema.BoolAttribute{
				Description: "Auto register environment files",
				Default:     booldefault.StaticBool(false),
				Required:    false,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *TorqueSpaceAdoServerRepositoryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueSpaceAdoServerRepositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueSpaceAdoServerRepositoryResourceModel
	const (
		StatusSyncing   = "Syncing"
		StatusConnected = "Connected"
		Interval        = 4 * time.Second
	)
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
	start := time.Now()
	err := r.client.OnboardAdoServerRepoToSpace(data.SpaceName.ValueString(), data.RepositoryName.ValueString(),
		data.RepositoryUrl.ValueString(), data.Token.ValueStringPointer(), data.Branch.ValueString(), data.CredentialName.ValueString(), agents, data.UseAllAgents.ValueBool(), data.AutoRegisterEac.ValueBool())
	if err != nil {
		repo, err := r.client.GetRepoDetails(data.SpaceName.ValueString(), data.RepositoryName.ValueString())
		if repo == nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to onboard repository to space, got error: %s", err))
			return
		}
		if repo.Status == StatusSyncing {
			timeout := time.Duration(data.TimeOut.ValueInt32()) * time.Minute
			for time.Since(start) < timeout {
				repo, err := r.client.GetRepoDetails(data.SpaceName.ValueString(), data.RepositoryName.ValueString())
				if err != nil {
					resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error while polling repository status: %s", err))
					return
				}
				if repo.Status == StatusConnected {
					resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
					return
				}
				time.Sleep(Interval)
			}
			resp.Diagnostics.AddError("Sync Timeout", "Timed out while syncing repository")
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to onboard repository to space, got error: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueSpaceAdoServerRepositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueSpaceAdoServerRepositoryResourceModel

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

func (r *TorqueSpaceAdoServerRepositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueSpaceAdoServerRepositoryResourceModel

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

func (r *TorqueSpaceAdoServerRepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueSpaceAdoServerRepositoryResourceModel

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

func (r *TorqueSpaceAdoServerRepositoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
