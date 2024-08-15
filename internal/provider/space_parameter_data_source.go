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
	_ datasource.DataSource              = &spaceParameterDataSource{}
	_ datasource.DataSourceWithConfigure = &spaceParameterDataSource{}
)

// NewspaceParametersDataSource is a helper function to simplify the provider implementation.
func NewSpaceParameterDataSource() datasource.DataSource {
	return &spaceParameterDataSource{}
}

// spaceParameterDataSource is the data source implementation.
type spaceParameterDataSource struct {
	client *client.Client
}

// spaceParameterDataSourceModel maps the data source schema data.
type spaceParameterDataSourceModel struct {
	SpaceName   types.String `tfsdk:"space_name"`
	Name        types.String `tfsdk:"name"`
	Value       types.String `tfsdk:"value"`
	Sensitive   types.Bool   `tfsdk:"sensitive"`
	Description types.String `tfsdk:"description"`
}

// Metadata returns the data source type name.
func (d *spaceParameterDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_space_parameter"
}

// Schema defines the schema for the data source.
func (d *spaceParameterDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get Account Parameter details.",
		Attributes: map[string]schema.Attribute{
			"space_name": schema.StringAttribute{
				MarkdownDescription: "The name of the Torque Space of the parameter",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the space level parameter",
				Required:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "Parameter Value. Value of sensitive parameter is null.",
				Computed:            true,
			},
			"sensitive": schema.BoolAttribute{
				MarkdownDescription: "Whether the parameter is sensitive or not.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Parameter description",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *spaceParameterDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *spaceParameterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state spaceParameterDataSourceModel
	var space_name types.String
	var parameter types.String

	diags := req.Config.GetAttribute(ctx, path.Root("space_name"), &space_name)
	diags = append(diags, req.Config.GetAttribute(ctx, path.Root("name"), &parameter)...)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	space_parameter_data, err := d.client.GetSpaceParameter(space_name.ValueString(), parameter.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Torque Account Parameter",
			err.Error(),
		)
		return
	}

	state.Name = types.StringValue(space_parameter_data.Name)
	if space_parameter_data.Sensitive {
		state.Value = types.StringNull()
	} else {
		state.Value = types.StringValue(space_parameter_data.Value)
	}
	state.SpaceName = space_name
	state.Sensitive = types.BoolValue(space_parameter_data.Sensitive)
	state.Description = types.StringValue(space_parameter_data.Description)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
