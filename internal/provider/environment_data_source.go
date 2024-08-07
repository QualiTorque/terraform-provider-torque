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
	_ datasource.DataSource              = &environmentDataSource{}
	_ datasource.DataSourceWithConfigure = &environmentDataSource{}
)

// NewenvironmentsDataSource is a helper function to simplify the provider implementation.
func NewEnvironmentDataSource() datasource.DataSource {
	return &environmentDataSource{}
}

// environmentDataSource is the data source implementation.
type environmentDataSource struct {
	client *client.Client
}

// environmentDataSourceModel maps the data source schema data.
type environmentDataSourceModel struct {
	SpaceName     types.String `tfsdk:"space_name"`
	Id            types.String `tfsdk:"id"`
	IsEAC         types.Bool   `tfsdk:"is_eac"`
	LastUsed      types.String `tfsdk:"last_used"`
	BlueprintName types.String `tfsdk:"blueprint_name"`
	OwnerEmail    types.String `tfsdk:"owner_email"`
}

type EnvironmentOwnerModel struct {
	OwnerEmail types.String `tfsdk:"email"`
}

type EnvironmentDetailsModel struct {
	Id             types.String               `tfsdk:"id"`
	ComputedStatus types.String               `tfsdk:"computed_status"`
	Definition     EnvironmentDefinitionModel `tfsdk:"definition"`
}

type EnvironmentDefinitionModel struct {
	Metadata EnvironmentMetadataModel `tfsdk:"metadata"`
}

type EnvironmentMetadataModel struct {
	BlueprintName types.String `tfsdk:"blueprint_name"`
}

type EnvironmentDetailTagsModel struct {
	Tags []TagModel `tfsdk:"tags"`
}

type TagModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type CollaboratorsModel struct {
	CollaboratorsEmails types.List `tfsdk:"collaborators_emails"`
	AllSpaceMembers     types.Bool `tfsdk:"all_space_members"`
}

// Metadata returns the data source type name.
func (d *environmentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

// Schema defines the schema for the data source.
func (d *environmentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get environment details.",
		Attributes: map[string]schema.Attribute{
			"space_name": schema.StringAttribute{
				MarkdownDescription: "Space",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Environment ID",
				Required:            true,
			},
			"blueprint_name": schema.StringAttribute{
				MarkdownDescription: "Blueprint Name",
				Computed:            true,
			},
			"is_eac": schema.BoolAttribute{
				MarkdownDescription: "Is environment source is Env-as-Code",
				Computed:            true,
			},
			"last_used": schema.StringAttribute{
				MarkdownDescription: "Last time environment was used",
				Computed:            true,
			},
			"owner_email": schema.StringAttribute{
				MarkdownDescription: "Last time environment was used",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *environmentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *environmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state environmentDataSourceModel
	var space_name types.String
	var id types.String

	// diags1 := req.Config.GetAttribute(ctx, path.Root("space_name"), &space_name)
	// resp.Diagnostics.Append(diags1...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// diags2 := req.Config.GetAttribute(ctx, path.Root("environment_id"), &environment_id)
	// resp.Diagnostics.Append(diags2...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// var data environmentDataSourceModel
	// var space_name types.String

	// diags := req.State.Get(ctx, &data)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	diags := req.Config.GetAttribute(ctx, path.Root("space_name"), &space_name)
	diags = req.Config.GetAttribute(ctx, path.Root("id"), &id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	environment_data, err := d.client.GetEnvironmentDetails(space_name.ValueString(), id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Torque environment",
			err.Error(),
		)
		return
	}

	state.LastUsed = types.StringValue(environment_data.LastUsed)
	state.IsEAC = types.BoolValue(environment_data.IsEAC)
	state.BlueprintName = types.StringValue(environment_data.Details.Definition.Metadata.BlueprintName)
	state.OwnerEmail = types.StringValue(environment_data.Owner.OwnerEmail)
	state.Id = types.StringValue(environment_data.Details.Id)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
