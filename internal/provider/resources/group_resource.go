package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/qualitorque/terraform-provider-torque/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TorqueGroupResource{}
var _ resource.ResourceWithImportState = &TorqueGroupResource{}

func NewTorqueGroupResource() resource.Resource {
	return &TorqueGroupResource{}
}

// TorqueGroupResource defines the resource implementation.
type TorqueGroupResource struct {
	client *client.Client
}

// TorqueGroupResourceModel describes the resource data model.
type TorqueGroupResourceModel struct {
	Name        types.String     `tfsdk:"group_name"`
	Description types.String     `tfsdk:"description"`
	IdpId       types.String     `tfsdk:"idp_identifier"`
	Users       types.List       `tfsdk:"users"`
	AccountRole types.String     `tfsdk:"account_role"`
	SpaceRoles  []SpaceRoleModel `tfsdk:"space_roles"`
}

type SpaceRoleModel struct {
	SpaceName types.String `tfsdk:"space_name"`
	SpaceRole types.String `tfsdk:"space_role"`
}

func (r *TorqueGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_group"
}

func (r *TorqueGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Creation of a new Torque group",

		Attributes: map[string]schema.Attribute{
			"group_name": schema.StringAttribute{
				MarkdownDescription: "The group name to create",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Group description to be presented in the Torque user interface",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"idp_identifier": schema.StringAttribute{
				MarkdownDescription: "Group association to IDP",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"users": schema.ListAttribute{
				MarkdownDescription: "Users to include in the newly created group",
				Optional:            true,
				Computed:            false,
				ElementType:         types.StringType,
			},
			"account_role": schema.StringAttribute{
				MarkdownDescription: "In case the group should be configured in the account level, use this attribute to define the group role in the account",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"space_roles": schema.ListNestedAttribute{
				Description: "key-value pairs of spaces and roles that the newly created group will be associated to",
				Computed:    false,
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"space_name": schema.StringAttribute{
							Description: "An existing Torque space name",
							Computed:    false,
							Optional:    true,
						},
						"space_role": schema.StringAttribute{
							Description: "Space role to be used for the specific space in the group",
							Computed:    false,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *TorqueGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var users []string
	if !data.Users.IsNull() {
		for _, user := range data.Users.Elements() {
			users = append(users, user.String())
		}
	}

	var spaceRoles []client.SpaceRole
	if len(data.SpaceRoles) > 0 {
		for _, spaceRole := range data.SpaceRoles {
			spaceRoles = append(spaceRoles, client.SpaceRole{
				SpaceName: spaceRole.SpaceName.ValueString(),
				SpaceRole: spaceRole.SpaceRole.ValueString(),
			})
		}
	}

	err := r.client.AddGroupToSpace(data.Name.ValueString(), data.Description.ValueString(), data.IdpId.ValueString(),
		users, data.AccountRole.ValueString(), spaceRoles)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create group, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Resource Created Successful!")

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *TorqueGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueGroupResourceModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.client.GetGroup(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading group details",
			"Could not read Torque group name "+data.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	// Treat HTTP 404 Not Found status as a signal to recreate resource
	// and return early
	if group.Name == "" {
		tflog.Error(ctx, "Group not found in Torque")
		resp.State.RemoveResource(ctx)
		return
	}

	data.AccountRole = types.StringValue(group.AccountRole)
	data.Description = types.StringValue(group.Description)
	data.IdpId = types.StringValue(group.IdpId)

	roles := []SpaceRoleModel{}
	for _, role := range group.SpaceRoles {
		roles = append(roles, SpaceRoleModel{
			SpaceName: types.StringValue(role.SpaceName),
			SpaceRole: types.StringValue(role.SpaceRole),
		})
	}

	data.SpaceRoles = roles

	data.Users, _ = types.ListValueFrom(ctx, types.StringType, group.Users)

	// Set refreshed state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *TorqueGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var data TorqueGroupResourceModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var roles []client.SpaceRole
	for _, role := range data.SpaceRoles {
		roles = append(roles, client.SpaceRole{
			SpaceName: role.SpaceName.ValueString(),
			SpaceRole: role.SpaceRole.ValueString(),
		})
	}

	var users []string
	if !data.Users.IsNull() {
		for _, user := range data.Users.Elements() {
			users = append(users, user.String())
		}
	}

	// Update existing order
	err := r.client.UpdateGroup(data.Name.ValueString(), data.Description.ValueString(), data.IdpId.ValueString(),
		users, data.AccountRole.ValueString(), roles)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Torque group",
			"Could not update group, unexpected error: "+err.Error(),
		)
		return
	}

	group, err := r.client.GetGroup(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading group details",
			"Could not read Torque group name "+data.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	data.AccountRole = types.StringValue(group.AccountRole)
	data.Description = types.StringValue(group.Description)
	data.IdpId = types.StringValue(group.IdpId)

	n_roles := []SpaceRoleModel{}
	for _, role := range group.SpaceRoles {
		n_roles = append(n_roles, SpaceRoleModel{
			SpaceName: types.StringValue(role.SpaceName),
			SpaceRole: types.StringValue(role.SpaceRole),
		})
	}
	data.SpaceRoles = n_roles

	data.Users, _ = types.ListValueFrom(ctx, types.StringType, group.Users)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TorqueGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueGroupResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteGroup(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete group, got error: %s", err))
		return
	}

}

func (r *TorqueGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("group_name"), req, resp)
}
