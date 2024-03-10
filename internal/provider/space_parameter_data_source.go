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
	_ datasource.DataSource              = &TorqueSpaceParameterDataSource{}
	_ datasource.DataSourceWithConfigure = &TorqueSpaceParameterDataSource{}
)

// NewusersDataSource is a helper function to simplify the provider implementation.
func NewTorqueSpaceParameterDataSource() datasource.DataSource {
	return &TorqueSpaceParameterDataSource{}
}

// userDataSource is the data source implementation.
type TorqueSpaceParameterDataSource struct {
	client *client.Client
}

// userDataSourceModel maps the data source schema data.
type TorqueSpaceParameterDataSourceModel struct {
	SpaceName   types.String `tfsdk:"space_name"`
	Name        types.String `tfsdk:"name"`
	Value       types.String `tfsdk:"value"`
	Sensitive   types.Bool   `tfsdk:"sensitive"`
	Description types.String `tfsdk:"description"`
}

// Metadata returns the data source type name.
func (d *TorqueSpaceParameterDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_space_parameter"
}

// Schema defines the schema for the data source.
func (d *TorqueSpaceParameterDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get blueprint information for a specific repository in a space",
		Attributes: map[string]schema.Attribute{
			"space_name": schema.StringAttribute{
				MarkdownDescription: "Space name to add the parameter to",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the new parameter to be added to torque",
				Required:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "Tag value to be set as the parameter in the space",
				Optional:            true,
				Computed:            false,
			},
			"sensitive": schema.BoolAttribute{
				MarkdownDescription: "Sensitive or not",
				Optional:            true,
				Computed:            false,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Parameter description",
				Optional:            true,
				Computed:            false,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *TorqueSpaceParameterDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *TorqueSpaceParameterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state TorqueSpaceParameterDataSourceModel
	var space_name types.String
	var name types.String

	diags1 := req.Config.GetAttribute(ctx, path.Root("space_name"), &space_name)
	resp.Diagnostics.Append(diags1...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags2 := req.Config.GetAttribute(ctx, path.Root("name"), &name)
	resp.Diagnostics.Append(diags2...)
	if resp.Diagnostics.HasError() {
		return
	}

	parameter_data, err := d.client.GetSpaceParameter(space_name.ValueString(), name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Space Parameter",
			err.Error(),
		)
		return
	}

	// initialize state
	state.Name = types.StringValue(parameter_data.Name)
	state.Value = types.StringValue(parameter_data.Value)
	state.Sensitive = types.BoolValue(parameter_data.Sensitive)
	state.Description = types.StringValue(parameter_data.Description)
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
