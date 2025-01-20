package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/qualitorque/terraform-provider-torque/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TorqueDeploymentEngineResource{}
var _ resource.ResourceWithImportState = &TorqueDeploymentEngineResource{}

func NewTorqueDeploymentEngineResource() resource.Resource {
	return &TorqueDeploymentEngineResource{}
}

// TorqueDeploymentEngineResource defines the resource implementation.
type TorqueDeploymentEngineResource struct {
	client *client.Client
}

// TorqueDeploymentEngineResourceModel describes the resource data model.
type TorqueDeploymentEngineResourceModel struct {
	Name                   types.String `tfsdk:"name"`
	Description            types.String `tfsdk:"description"`
	Type                   types.String `tfsdk:"type"`
	AuthToken              types.String `tfsdk:"auth_token"`
	AgentName              types.String `tfsdk:"agent_name"`
	ServerUrl              types.String `tfsdk:"server_url"`
	PollingIntervalSeconds types.Int32  `tfsdk:"polling_interval_seconds"`
	AllSpaces              types.Bool   `tfsdk:"all_spaces"`
	SpecificSpaces         types.List   `tfsdk:"specific_spaces"`
}

const (
	argocd_engine_type = "Argo CD"
)

func (r *TorqueDeploymentEngineResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_deployment_engine"
}

func (r *TorqueDeploymentEngineResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Creates a new deployment engine.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the deployment engine.",
				Optional:            false,
				Computed:            false,
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the deployment engine",
				Optional:            true,
				Computed:            true,
				Required:            false,
				Default:             stringdefault.StaticString(""),
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of the deployment engine.",
				Optional:            false,
				Computed:            true,
				Required:            false,
				Default:             stringdefault.StaticString("Argo CD"),
			},
			"agent_name": schema.StringAttribute{
				MarkdownDescription: "Type of the deployment engine.",
				Optional:            false,
				Computed:            false,
				Required:            true,
			},
			"server_url": schema.StringAttribute{
				MarkdownDescription: "Server URL of the deployment engine",
				Optional:            false,
				Computed:            false,
				Required:            true,
			},
			"auth_token": schema.StringAttribute{
				MarkdownDescription: "Token of the deployment engine.",
				Optional:            false,
				Required:            true,
				Sensitive:           true,
			},
			"polling_interval_seconds": schema.Int32Attribute{
				MarkdownDescription: "Polling interval of the deployment engine in seconds.",
				Optional:            true,
				Computed:            true,
				Required:            false,
				Default:             int32default.StaticInt32(30),
			},
			"all_spaces": schema.BoolAttribute{
				Description: "Specify if the deployment engine can be used in all spaces. Defaults to true, use specific spaces attribute for allowing only specific spaces.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					&allSpacesModifier{},
				},
				Validators: []validator.Bool{
					// Validate only this attribute or other_attr is configured or neither.
					boolvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("specific_spaces"),
					}...),
				},
			},
			"specific_spaces": schema.ListAttribute{
				Description: "List of spaces that can use this deployment engine",
				Required:    false,
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					// Validate only this attribute or other_attr is configured or neither.
					listvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("all_spaces"),
					}...),
				},
			},
		},
	}
}
func (r *TorqueDeploymentEngineResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueDeploymentEngineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueDeploymentEngineResourceModel
	var allowed_spaces client.AllowedSpaces
	var specificSpaces []string
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	allowed_spaces.AllSpaces = data.AllSpaces.ValueBool()
	if !data.SpecificSpaces.IsNull() {
		allowed_spaces.AllSpaces = false
		for _, val := range data.SpecificSpaces.Elements() {
			specificSpaces = append(specificSpaces, strings.Replace(val.String(), "\"", "", -1))
		}
		allowed_spaces.SpecificSpaces = specificSpaces
	} else {
		allowed_spaces.AllSpaces = data.AllSpaces.ValueBool() // true
	}

	err := r.client.CreateDeploymentEngine(argocd_engine_type, data.Name.ValueString(), data.Description.ValueString(), data.AgentName.ValueString(), data.AuthToken.ValueString(), data.PollingIntervalSeconds.ValueInt32(), data.ServerUrl.ValueString(), allowed_spaces)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create deployment engine, got error: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueDeploymentEngineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueDeploymentEngineResourceModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deployment_engine, err := r.client.GetDeploymentEngine(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading deployment engine details",
			"Could not read deployment engine "+data.Name.ValueString()+": "+err.Error(),
		)
		return
	}
	data.Name = types.StringValue(deployment_engine.Name)
	data.Description = types.StringValue(deployment_engine.Description)
	data.Type = types.StringValue(deployment_engine.Type)
	data.AgentName = types.StringValue(deployment_engine.Agent.Name)
	data.ServerUrl = types.StringValue(deployment_engine.ServerUrl)
	data.PollingIntervalSeconds = types.Int32Value(int32(deployment_engine.PollingIntervalSeconds.Value))

	data.AllSpaces = types.BoolValue(deployment_engine.AllowedSpaces.AllSpaces)
	if len(deployment_engine.AllowedSpaces.SpecificSpaces) > 0 {
		data.SpecificSpaces, _ = types.ListValueFrom(ctx, types.StringType, deployment_engine.AllowedSpaces.SpecificSpaces)
	} else {
		data.SpecificSpaces = types.ListNull(types.StringType)
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TorqueDeploymentEngineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueDeploymentEngineResourceModel
	var state TorqueDeploymentEngineResourceModel

	var allowed_spaces client.AllowedSpaces
	var specificSpaces []string

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	allowed_spaces.AllSpaces = data.AllSpaces.ValueBool()
	if !data.SpecificSpaces.IsNull() {
		allowed_spaces.AllSpaces = false
		for _, val := range data.SpecificSpaces.Elements() {
			specificSpaces = append(specificSpaces, strings.Replace(val.String(), "\"", "", -1))
		}
		allowed_spaces.SpecificSpaces = specificSpaces
	} else {
		allowed_spaces.AllSpaces = data.AllSpaces.ValueBool() // true
	}

	err := r.client.UpdateDeploymentEngine(argocd_engine_type, state.Name.ValueString(), data.Name.ValueString(), data.Description.ValueString(), data.AgentName.ValueString(), data.AuthToken.ValueString(), data.PollingIntervalSeconds.ValueInt32(), data.ServerUrl.ValueString(), allowed_spaces)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update deployment engine, got error: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueDeploymentEngineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueDeploymentEngineResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDeploymentEngine(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete deployment engine, got error: %s", err))
		return
	}

}

func (r *TorqueDeploymentEngineResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
