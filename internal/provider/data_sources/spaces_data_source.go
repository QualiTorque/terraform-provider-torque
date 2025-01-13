package data_sources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/qualitorque/terraform-provider-torque/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &SpacesDataSource{}
	_ datasource.DataSourceWithConfigure = &SpacesDataSource{}
)

// NewSpacessDataSource is a helper function to simplify the provider implementation.
func NewSpacesDataSource() datasource.DataSource {
	return &SpacesDataSource{}
}

// SpacesDataSource is the data source implementation.
type SpacesDataSource struct {
	client *client.Client
}

// SpacesDataSourceModel maps the data source schema data.
type SpacesDataSourceModel struct {
	Spaces []spaceModel `tfsdk:"spaces"`
}

type spaceModel struct {
	Name        types.String `tfsdk:"name"`
	NumOfUsers  types.Int32  `tfsdk:"num_of_users"`
	NumOfGroups types.Int32  `tfsdk:"num_of_groups"`
	Color       types.String `tfsdk:"color"`
	Icon        types.String `tfsdk:"icon"`
}

// Metadata returns the data source type name.
func (d *SpacesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_spaces"
}

// Schema defines the schema for the data source.
func (d *SpacesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a list of all spaces in the Torque account..",
		Attributes: map[string]schema.Attribute{
			"spaces": schema.ListNestedAttribute{
				Description: "Spaces in the account",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the space",
							Computed:            true,
						},
						"num_of_users": schema.Int32Attribute{
							MarkdownDescription: "Number of users in the space",
							Computed:            true,
						},
						"num_of_groups": schema.Int32Attribute{
							MarkdownDescription: "Number of groups in the space",
							Computed:            true,
						},
						"color": schema.StringAttribute{
							MarkdownDescription: "Space color",
							Computed:            true,
						},
						"icon": schema.StringAttribute{
							MarkdownDescription: "Space Icon",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *SpacesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *SpacesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state SpacesDataSourceModel

	spaces, err := d.client.GetSpaces()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Torque Spaces",
			err.Error(),
		)
		return
	}
	for _, space_item := range spaces {
		space := spaceModel{
			Name:        types.StringValue(space_item.Name),
			Color:       types.StringValue(space_item.Color),
			Icon:        types.StringValue(space_item.Icon),
			NumOfUsers:  types.Int32Value(space_item.NumOfUsers),
			NumOfGroups: types.Int32Value(space_item.NumOfGroups),
		}
		state.Spaces = append(state.Spaces, space)
	}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
