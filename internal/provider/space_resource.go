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
var _ resource.Resource = &TorqueSpaceResource{}
var _ resource.ResourceWithImportState = &TorqueSpaceResource{}

func NewTorqueSpaceResource() resource.Resource {
	return &TorqueSpaceResource{}
}

// TorqueSpaceResource defines the resource implementation.
type TorqueSpaceResource struct {
	client *client.Client
}

// TorqueSpaceResourceModel describes the resource data model.
type TorqueSpaceResourceModel struct {
	Name                       types.String           `tfsdk:"name"`
	Color                      types.String           `tfsdk:"color"`
	Icon                       types.String           `tfsdk:"icon"`
	AssociatedMembers          types.List             `tfsdk:"space_members"`
	AssociatedAdmins           types.List             `tfsdk:"space_admins"`
	AssociatedKubernetesAgents []KubAgentRequestModel `tfsdk:"associated_kubernetes_agent"`
	AssociatedRepos            []RepoRequestModel     `tfsdk:"associated_repos"`
}

type KubAgentRequestModel struct {
	AgentName             types.String `tfsdk:"agent_name"`
	DefaultNamespace      types.String `tfsdk:"default_namespace"`
	DefaultServiceAccount types.String `tfsdk:"default_service_account"`
}

type RepoRequestModel struct {
	RepoName   types.String `tfsdk:"repository_name"`
	RepoBranch types.String `tfsdk:"branch"`
	RepoType   types.String `tfsdk:"repository_type"`
	RepoToken  types.String `tfsdk:"access_token"`
	RepoUrl    types.String `tfsdk:"repository_url"`
}

func (r *TorqueSpaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_space_resource"
}

func (r *TorqueSpaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Creation of a new Torque space with associated entities (users, repos, etc...)",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Space name to be create",
				Required:            true,
			},
			"color": schema.StringAttribute{
				MarkdownDescription: "Space color to be used for the new space",
				Optional:            true,
				Computed:            false,
			},
			"icon": schema.StringAttribute{
				MarkdownDescription: "Space icon to be used",
				Optional:            true,
				Computed:            false,
			},
			"space_members": schema.ListAttribute{
				MarkdownDescription: "List of space members to be associate to the newly created space",
				Optional:            true,
				Computed:            false,
				ElementType:         types.StringType,
			},
			"space_admins": schema.ListAttribute{
				MarkdownDescription: "List of space admins to be associate to the newly created space",
				Optional:            true,
				Computed:            false,
				ElementType:         types.StringType,
			},
			"associated_kubernetes_agent": schema.ListNestedAttribute{
				MarkdownDescription: "Kubernetes agent to associate to the newly create space",
				Optional:            true,
				Computed:            false,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"agent_name": schema.StringAttribute{
							Description: "Agent name to associate to the newly created space",
							Required:    true,
						},
						"default_service_account": schema.StringAttribute{
							Description: "Default service account to be used with the agent in the space",
							Required:    true,
						},
						"default_namespace": schema.StringAttribute{
							Description: "Default namespace to be used with the agent in the space",
							Required:    true,
						},
					},
				},
			},
			"associated_repos": schema.ListNestedAttribute{
				MarkdownDescription: "Kubernetes agent to associate to the newly create space",
				Optional:            true,
				Computed:            false,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"repository_name": schema.StringAttribute{
							Description: "The name of the repository to onboard in the newly created space",
							Required:    true,
						},
						"branch": schema.StringAttribute{
							Description: "Repository branch to use for blueprints and automation assets",
							Optional:    true,
						},
						"repository_type": schema.StringAttribute{
							Description: "Repository type. Available types: github, bitbucket, gitlab, azure (for Azure DevOps)",
							Required:    true,
						},
						"access_token": schema.StringAttribute{
							Description: "Personal Access Token (PAT) to authenticate with to the repository",
							Required:    true,
						},
						"repository_url": schema.StringAttribute{
							Description: "Repository URL. For example: https://github.com/<org>/<repo>",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func (r *TorqueSpaceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func trimQuote(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}

func (r *TorqueSpaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueSpaceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CreateSpace(data.Name.ValueString(), data.Color.ValueString(), data.Icon.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create space, got error: %s", err))
		return
	}

	if !data.AssociatedMembers.IsNull() {
		for _, member := range data.AssociatedMembers.Elements() {
			err := r.client.AddUserToSpace(trimQuote(member.String()), "Space Member", data.Name.ValueString())
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to attach space member to space, got error: %s", err))
				return
			}
		}
	}

	if !data.AssociatedAdmins.IsNull() {
		for _, admin := range data.AssociatedAdmins.Elements() {
			err := r.client.AddUserToSpace(trimQuote(admin.String()), "Space Admin", data.Name.ValueString())
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to attach space admin to space, got error: %s", err))
				return
			}
		}
	}

	for _, associationRequest := range data.AssociatedKubernetesAgents {
		err := r.client.AddAgentToSpace(associationRequest.AgentName.ValueString(), associationRequest.DefaultNamespace.ValueString(),
			associationRequest.DefaultServiceAccount.ValueString(), data.Name.ValueString(), "K8S")
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to attach agent to space, got error: %s", err))
			return
		}
	}

	for _, associationRepoRequest := range data.AssociatedRepos {
		err := r.client.OnboardRepoToSpace(data.Name.ValueString(), associationRepoRequest.RepoName.ValueString(), trimQuote(associationRepoRequest.RepoType.String()),
			associationRepoRequest.RepoUrl.ValueString(), associationRepoRequest.RepoToken.ValueString(),
			associationRepoRequest.RepoBranch.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to onboard repository to space, got error: %s", err))
			return
		}
	}

	tflog.Trace(ctx, "Resource Created Successful!")

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueSpaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueSpaceResourceModel

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

func (r *TorqueSpaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueSpaceResourceModel

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

func (r *TorqueSpaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueSpaceResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Remove repos from space.
	for _, associationRequest := range data.AssociatedKubernetesAgents {
		err := r.client.RemoveAgentFromSpace(associationRequest.AgentName.ValueString(), data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to attach agent to space, got error: %s", err))
			return
		}
	}

	// Delete the space.
	err := r.client.DeleteSpace(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete space, got error: %s", err))
		return
	}

}

func (r *TorqueSpaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
