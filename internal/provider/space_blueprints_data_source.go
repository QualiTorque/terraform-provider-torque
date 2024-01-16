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
	_ datasource.DataSource              = &spaceRepoBlueprintsDataSource{}
	_ datasource.DataSourceWithConfigure = &spaceRepoBlueprintsDataSource{}
)

// NewusersDataSource is a helper function to simplify the provider implementation.
func NewSpaceRepositoryBlueprintsDataSource() datasource.DataSource {
	return &spaceRepoBlueprintsDataSource{}
}

// userDataSource is the data source implementation.
type spaceRepoBlueprintsDataSource struct {
	client *client.Client
}

// userDataSourceModel maps the data source schema data.
type spaceRepoDataSourceModel struct {
	SpaceName  types.String     `tfsdk:"space_name"`
	Blueprints []blueprintModel `tfsdk:"blueprints"`
}

type blueprintModel struct {
	BlueprintName types.String `tfsdk:"blueprint_name"`
	Name          types.String `tfsdk:"name"`
	DisplayName   types.String `tfsdk:"display_name"`
	RepoName      types.String `tfsdk:"repository_name"`
	RepoBranch    types.String `tfsdk:"repository_branch"`
	Commit        types.String `tfsdk:"commit"`
	Description   types.String `tfsdk:"description"`
	Url           types.String `tfsdk:"url"`
	ModifiedBy    types.String `tfsdk:"modified_by"`
	Published     types.Bool   `tfsdk:"enabled"`
}

// Metadata returns the data source type name.
func (d *spaceRepoBlueprintsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_space_blueprints"
}

// Schema defines the schema for the data source.
func (d *spaceRepoBlueprintsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get blueprint information for a specific repository in a space",
		Attributes: map[string]schema.Attribute{
			"space_name": schema.StringAttribute{
				MarkdownDescription: "The name of the space to use",
				Required:            true,
			},
			"blueprints": schema.ListNestedAttribute{
				Description: "Blueprints in the space",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"blueprint_name": schema.StringAttribute{
							Description: "The unique name of the blueprint",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "no idea!! ",
							Computed:    true,
						},
						"display_name": schema.StringAttribute{
							Description: "The user friendly name of the blueprint in the space",
							Computed:    true,
						},
						"repository_name": schema.StringAttribute{
							Description: "The repository name from which the blueprint is used",
							Computed:    true,
						},
						"repository_branch": schema.StringAttribute{
							Description: "The branch from which the blueprint is taken",
							Computed:    true,
						},
						"commit": schema.StringAttribute{
							Description: "The commit id of the blueprint",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "The description of the bluepring",
							Computed:    true,
						},
						"url": schema.StringAttribute{
							Description: "URI of the blueprint",
							Computed:    true,
						},
						"modified_by": schema.StringAttribute{
							Description: "The name of the user that last modified the blueprint",
							Computed:    true,
						},
						"enabled": schema.BoolAttribute{
							Description: "is Published blueprint in the space",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *spaceRepoBlueprintsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *spaceRepoBlueprintsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state spaceRepoDataSourceModel

	var space types.String

	diags := req.Config.GetAttribute(ctx, path.Root("space_name"), &space)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	blueprints_data, err := d.client.GetSpaceBlueprints(space.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Torque user",
			err.Error(),
		)
		return
	}

	state.SpaceName = types.StringValue(space.ValueString())

	for _, blueprintItem := range blueprints_data {
		blueprintData := blueprintModel{
			BlueprintName: types.StringValue(blueprintItem.BlueprintName),
			Name:          types.StringValue(blueprintItem.Name),
			RepoName:      types.StringValue(blueprintItem.RepoName),
			Description:   types.StringValue(blueprintItem.Description),
			Commit:        types.StringValue(blueprintItem.Commit),
			ModifiedBy:    types.StringValue(blueprintItem.ModifiedBy),
			DisplayName:   types.StringValue(blueprintItem.DisplayName),
			RepoBranch:    types.StringValue(blueprintItem.RepoBranch),
			Url:           types.StringValue(blueprintItem.Url),
			Published:     types.BoolValue(blueprintItem.Published),
		}
		state.Blueprints = append(state.Blueprints, blueprintData)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
