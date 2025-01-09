package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
var _ resource.Resource = &TorqueGitCredentialsResource{}
var _ resource.ResourceWithImportState = &TorqueGitCredentialsResource{}

func NewTorqueGitCredentialsResource() resource.Resource {
	return &TorqueGitCredentialsResource{}
}

// TorqueGitCredentialsResource defines the resource implementation.
type TorqueGitCredentialsResource struct {
	client *client.Client
}

// TorqueGitCredentialsResourceModel describes the resource data model.
type TorqueGitCredentialsResourceModel struct {
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	Token             types.String `tfsdk:"token"`
	Type              types.String `tfsdk:"type"`
	AllowedSpaceNames types.List   `tfsdk:"allowed_space_names"`
	CloudType         types.String `tfsdk:"cloudtype"`
}

func (r *TorqueGitCredentialsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_git_credentials"
}

func (r *TorqueGitCredentialsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creation of a new git credentials resource, which can later be used to onboard a git repository. Supported repositories are github, gitlab enterprise, azure devops and bitbucket.",

		Attributes: map[string]schema.Attribute{
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
			"cloudtype": schema.StringAttribute{
				MarkdownDescription: "Credentials type identifier.",
				Required:            false,
				Computed:            true,
				Default:             stringdefault.StaticString("sourceControl"),
			},
			"allowed_space_names": schema.ListAttribute{
				MarkdownDescription: "List of allowed spaces that can use the credentials. At least one space must be in the list if a list is provided. If the argument is not probvided, the credentials may be used in all spaces",
				Required:            false,
				Optional:            true,
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1), // Ensure the list has at least one entry if required
				},
			},
		},
	}
}

func (r *TorqueGitCredentialsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueGitCredentialsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueGitCredentialsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	allowed_space_names := []string{}
	if !data.AllowedSpaceNames.IsNull() {
		for _, name := range data.AllowedSpaceNames.Elements() {
			allowed_space_names = append(allowed_space_names, strings.Trim(name.String(), "\""))
		}
	}
	err := r.client.CreateAccountCredentials(data.Name.ValueString(), data.Description.ValueString(), data.CloudType.ValueString(), data.Type.ValueString(), data.Type.ValueString(), data.Token.ValueStringPointer(), nil, nil, allowed_space_names)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create account git credentials, got error: %s", err))
		return
	}
	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *TorqueGitCredentialsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueGitCredentialsResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TorqueGitCredentialsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueGitCredentialsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	allowed_space_names := []string{}
	if !data.AllowedSpaceNames.IsNull() {
		for _, name := range data.AllowedSpaceNames.Elements() {
			allowed_space_names = append(allowed_space_names, strings.Trim(name.String(), "\""))
		}
	}

	err := r.client.UpdateAccountCredentials(data.Name.ValueString(), data.Description.ValueString(), data.CloudType.ValueString(), data.Type.ValueString(), data.Type.ValueString(), data.Token.ValueStringPointer(), nil, nil, allowed_space_names)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update git credentials, got error: %s", err))
		return
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueGitCredentialsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueGitCredentialsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.DeleteAccountCredentials(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete git credentials, got error: %s", err))
		return
	}
}

func (r *TorqueGitCredentialsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
