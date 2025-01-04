package data_sources

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
	_ datasource.DataSource              = &spaceCustomIconDataSource{}
	_ datasource.DataSourceWithConfigure = &spaceCustomIconDataSource{}
)

// NewaccountParametersDataSource is a helper function to simplify the provider implementation.
func NewSpaceCustomIconDataSource() datasource.DataSource {
	return &spaceCustomIconDataSource{}
}

// spaceCustomIconDataSource is the data source implementation.
type spaceCustomIconDataSource struct {
	client *client.Client
}

// spaceCustomIconDataSourceModel maps the data source schema data.
type spaceCustomIconDataSourceModel struct {
	SpaceName types.String `tfsdk:"space_name"`
	FileName  types.String `tfsdk:"file_name"`
	Key       types.String `tfsdk:"key"`
}

// Metadata returns the data source type name.
func (d *spaceCustomIconDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_space_custom_icon"
}

// Schema defines the schema for the data source.
func (d *spaceCustomIconDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get Space Custom Icon details.",
		Attributes: map[string]schema.Attribute{
			"space_name": schema.StringAttribute{
				MarkdownDescription: "Name of the space where this custom icon exists",
				Required:            true,
			},
			"file_name": schema.StringAttribute{
				MarkdownDescription: "Icon SVG file name.",
				Required:            true,
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "Identifier for the icon, to be used in the catalog item resource when associating this icon with a catalog item.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *spaceCustomIconDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *spaceCustomIconDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state spaceCustomIconDataSourceModel
	var file_name types.String
	var space_name types.String

	diags := req.Config.GetAttribute(ctx, path.Root("space_name"), &space_name)
	diags = append(diags, req.Config.GetAttribute(ctx, path.Root("file_name"), &file_name)...)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	custom_icon, err := d.client.GetCustomIcon(space_name.ValueString(), file_name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Torque Custom Icon",
			err.Error(),
		)
		return
	}

	state.SpaceName = space_name
	state.FileName = types.StringValue(custom_icon.FileName)
	state.Key = types.StringValue(custom_icon.Key)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
