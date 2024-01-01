package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/qualitorque/terraform-provider-torque/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &userDataSource{}
	_ datasource.DataSourceWithConfigure = &userDataSource{}
)

// NewusersDataSource is a helper function to simplify the provider implementation.
func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

// userDataSource is the data source implementation.
type userDataSource struct {
	client *client.Client
}

// userDataSourceModel maps the data source schema data.
type userDataSourceModel struct {
	UserEmail            types.String `tfsdk:"user_email"`
	FirstName            types.String `tfsdk:"first_name"`
	LastName             types.String `tfsdk:"last_name"`
	Timezone             types.String `tfsdk:"timezone"`
	DisplayFirstName     types.String `tfsdk:"display_first_name"`
	DisplayLastName      types.String `tfsdk:"display_last_name"`
	UserType             types.String `tfsdk:"user_type"`
	JoinDate             types.String `tfsdk:"join_date"`
	AccountRole          types.String `tfsdk:"account_role"`
	HasAccessToAllSpaces types.Bool   `tfsdk:"has_access_to_all_spaces"`
	Permissions          types.List   `tfsdk:"permissions"`
}

// Metadata returns the data source type name.
func (d *userDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

// Schema defines the schema for the data source.
func (d *userDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get user details.",
		Attributes: map[string]schema.Attribute{
			"user_email": schema.StringAttribute{
				MarkdownDescription: "The Email of the user",
				Optional:            true,
			},
			"first_name": schema.StringAttribute{
				MarkdownDescription: "User first name",
				Computed:            true,
			},
			"last_name": schema.StringAttribute{
				MarkdownDescription: "User last name",
				Computed:            true,
			},
			"timezone": schema.StringAttribute{
				MarkdownDescription: "User timezone",
				Computed:            true,
			},
			"display_first_name": schema.StringAttribute{
				MarkdownDescription: "User display first name",
				Computed:            true,
			},
			"display_last_name": schema.StringAttribute{
				MarkdownDescription: "User display last name",
				Computed:            true,
			},
			"user_type": schema.StringAttribute{
				MarkdownDescription: "User type",
				Computed:            true,
			},
			"account_role": schema.StringAttribute{
				MarkdownDescription: "User role in the account",
				Computed:            true,
			},
			"join_date": schema.StringAttribute{
				MarkdownDescription: "User joined date to Torque",
				Computed:            true,
			},
			"has_access_to_all_spaces": schema.BoolAttribute{
				MarkdownDescription: "User access to all space",
				Computed:            true,
			},
			"permissions": schema.ListAttribute{
				MarkdownDescription: "User permissions list",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state userDataSourceModel

	var email types.String

	diags1 := req.Config.GetAttribute(ctx, path.Root("user_email"), &email)
	resp.Diagnostics.Append(diags1...)
	if resp.Diagnostics.HasError() {
		return
	}

	user_data, err := d.client.GetUserDetails(email.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Torque users",
			err.Error(),
		)
		return
	}

	state.UserEmail = types.StringValue(user_data.Email)
	state.FirstName = types.StringValue(user_data.FirstName)
	state.LastName = types.StringValue(user_data.LastName)
	state.AccountRole = types.StringValue(user_data.AccountRole)
	state.DisplayFirstName = types.StringValue(user_data.DisplayFirstName)
	state.DisplayLastName = types.StringValue(user_data.DisplayLastName)
	state.HasAccessToAllSpaces = types.BoolValue(user_data.HasAccessToAllSpaces)
	state.JoinDate = types.StringValue(user_data.JoinDate)
	state.Timezone = types.StringValue(user_data.Timezone)
	state.UserType = types.StringValue(user_data.UserType)
	state.Permissions, _ = types.ListValueFrom(ctx, types.StringType, user_data.Permissions)
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
