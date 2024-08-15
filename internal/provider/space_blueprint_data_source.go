package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/qualitorque/terraform-provider-torque/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &spaceBlueprintDataSource{}
	_ datasource.DataSourceWithConfigure = &spaceBlueprintDataSource{}
)

// NewusersDataSource is a helper function to simplify the provider implementation.
func NewSpaceBlueprintDataSource() datasource.DataSource {
	return &spaceBlueprintDataSource{}
}

// userDataSource is the data source implementation.
type spaceBlueprintDataSource struct {
	client *client.Client
}

// userDataSourceModel maps the data source schema data.
type spaceBlueprintDataSourceModel struct {
	SpaceName    types.String `tfsdk:"space_name"`
	Name         types.String `tfsdk:"name"`
	DisplayName  types.String `tfsdk:"display_name"`
	RepoName     types.String `tfsdk:"repository_name"`
	RepoBranch   types.String `tfsdk:"repository_branch"`
	Commit       types.String `tfsdk:"commit"`
	Description  types.String `tfsdk:"description"`
	Url          types.String `tfsdk:"url"`
	ModifiedBy   types.String `tfsdk:"modified_by"`
	LastModified types.String `tfsdk:"last_modified"`
	Published    types.Bool   `tfsdk:"enabled"`
	// Details   blueprint     `tfsdk:"details"`
	Tags                  []blueprintTag `tfsdk:"tags"`
	MaxDuration           types.String   `tfsdk:"max_duration"`
	DefaultDuration       types.String   `tfsdk:"default_duration"`
	DefaultExtend         types.String   `tfsdk:"default_extend"`
	MaxActiveEnvironments types.Int32    `tfsdk:"max_active_environments"`
	AlwaysOn              types.Bool     `tfsdk:"always_on"`
	// Policies  policies      `tfsdk:"policies"`
	Outputs types.List `tfsdk:"outputs"`
}

// type blueprintOutput struct {
// 	Name types.String `tfsdk:"name"`
// }

type blueprintTag struct {
	Name           types.String `tfsdk:"name"`
	DefaultValue   types.String `tfsdk:"default_value"`
	PossibleValues types.List   `tfsdk:"possible_values"`
	Description    types.String `tfsdk:"description"`
}

// type policies struct {
// 	MaxDuration           types.String `tfsdk:"max_duration"`
// 	DefaultDuration       types.String `tfsdk:"default_duration"`
// 	DefaultExtend         types.String `tfsdk:"default_extend"`
// 	MaxActiveEnvironments types.Int32  `tfsdk:"max_active_environments"`
// 	AlwaysOn              types.Bool   `tfsdk:"always_on"`
// }

// type blueprint struct {
// 	BlueprintName types.String `tfsdk:"blueprint_name"`
// 	Name          types.String `tfsdk:"name"`
// 	DisplayName   types.String `tfsdk:"display_name"`
// 	RepoName      types.String `tfsdk:"repository_name"`
// 	RepoBranch    types.String `tfsdk:"repository_branch"`
// 	Commit        types.String `tfsdk:"commit"`
// 	Description   types.String `tfsdk:"description"`
// 	Url           types.String `tfsdk:"url"`
// 	ModifiedBy    types.String `tfsdk:"modified_by"`
// 	Published     types.Bool   `tfsdk:"enabled"`
// }

// Metadata returns the data source type name.
func (d *spaceBlueprintDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_space_blueprint"
}

// Schema defines the schema for the data source.
func (d *spaceBlueprintDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get blueprint information for a specific repository in a space",
		Attributes: map[string]schema.Attribute{
			"space_name": schema.StringAttribute{
				MarkdownDescription: "The name of the space to use",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the blueprint",
				Required:            true,
			},
			"display_name": schema.StringAttribute{
				Description: "The user-friendly name of the blueprint in the space",
				Computed:    true,
			},
			"repository_name": schema.StringAttribute{
				Description: "The repository name from which the blueprint is used",
				Computed:    true,
				Optional:    true,
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
				Description: "The description of the blueprint",
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
			"last_modified": schema.StringAttribute{
				Description: "The name of the user that last modified the blueprint",
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Is the published blueprint in the space",
				Computed:    true,
			},
			"tags": schema.ListNestedAttribute{
				Description: "Blueprints in the space",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The unique name of the blueprint",
							Computed:    true,
						},
						"default_value": schema.StringAttribute{
							Description: "The blueprint display name ",
							Computed:    true,
						},
						"possible_values": schema.ListAttribute{
							Description: "The user friendly name of the blueprint in the space",
							Computed:    true,
							ElementType: types.StringType,
						},
						"description": schema.StringAttribute{
							Description: "The repository name from which the blueprint is used",
							Computed:    true,
						},
					},
				},
			},
			"outputs": schema.ListAttribute{
				Description: "Blueprints in the space",
				Computed:    true,
				ElementType: types.StringType,
			},
			// "outputs": schema.ListNestedAttribute{
			// 	Description: "Blueprints in the space",
			// 	Computed:    true,
			// 	NestedObject: schema.NestedAttributeObject{
			// 		Attributes: map[string]schema.Attribute{
			// 			"name": schema.StringAttribute{
			// 				Description: "The unique name of the blueprint",
			// 				Computed:    true,
			// 			},
			// 		},
			// 	},
			// },
			"max_duration": schema.StringAttribute{
				Description: "The unique name of the blueprint",
				Computed:    true,
			},
			"default_duration": schema.StringAttribute{
				Description: "The unique name of the blueprint",
				Computed:    true,
			},
			"default_extend": schema.StringAttribute{
				Description: "The unique name of the blueprint",
				Computed:    true,
			},
			"max_active_environments": schema.Int32Attribute{
				Description: "The unique name of the blueprint",
				Computed:    true,
			},
			"always_on": schema.BoolAttribute{
				Description: "The unique name of the blueprint",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *spaceBlueprintDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *spaceBlueprintDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state spaceBlueprintDataSourceModel
	var name types.String
	var space_name types.String

	diags := req.Config.GetAttribute(ctx, path.Root("space_name"), &space_name)
	diags = append(diags, req.Config.GetAttribute(ctx, path.Root("name"), &name)...)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	blueprint_data, err := d.client.GetBlueprintDetails(space_name.ValueString(), name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Torque user",
			err.Error(),
		)
		return
	}

	// initialize state
	state.SpaceName = types.StringValue(space_name.ValueString())
	state.Name = types.StringValue(name.ValueString())
	state.DisplayName = types.StringValue(blueprint_data.Details.DisplayName)
	state.Commit = types.StringValue(blueprint_data.Details.Commit)
	state.Description = types.StringValue(blueprint_data.Details.Description)
	state.Published = types.BoolValue(blueprint_data.Details.Published)
	state.MaxDuration = types.StringValue(blueprint_data.Policies.MaxDuration)
	state.DefaultDuration = types.StringValue(blueprint_data.Policies.DefaultDuration)
	state.DefaultExtend = types.StringValue(blueprint_data.Policies.DefaultExtend)
	state.MaxActiveEnvironments = types.Int32Value(blueprint_data.Policies.MaxActiveEnvironments)
	state.AlwaysOn = types.BoolValue(blueprint_data.Policies.AlwaysOn)
	state.ModifiedBy = types.StringValue(blueprint_data.Details.ModifiedBy)
	state.LastModified = types.StringValue(blueprint_data.Details.LastModified)
	state.RepoBranch = types.StringValue(blueprint_data.Details.RepoBranch)
	state.RepoName = types.StringValue(blueprint_data.Details.RepoName)
	state.Url = types.StringValue(blueprint_data.Details.Url)

	for _, tagItem := range blueprint_data.Tags {
		var possibleValues []attr.Value
		for _, value := range tagItem.PossibleValues {
			possibleValues = append(possibleValues, types.StringValue(value))
		}
		possibleValuesList, _ := types.ListValue(types.StringType, possibleValues)
		tagData := blueprintTag{
			Name:           types.StringValue(tagItem.Name),
			DefaultValue:   types.StringValue(tagItem.DefaultValue),
			PossibleValues: possibleValuesList,
			Description:    types.StringValue(tagItem.Description),
		}
		state.Tags = append(state.Tags, tagData)
	}

	var outputs []attr.Value
	for _, output := range blueprint_data.Details.Outputs {
		outputs = append(outputs, types.StringValue(output.Name))
	}
	outputsList, _ := types.ListValue(types.StringType, outputs)
	state.Outputs = outputsList
	// for _, outputItem := range blueprint_data.Details.Outputs {
	// 	// outputData := blueprintOutput{
	// 	// 	Name: types.StringValue(outputItem.Name),
	// 	// }
	// 	state.Outputs = append(state.Outputs, outputData)
	// }

	// if !state.RepoFilter.IsNull() {
	// 	for _, blueprintItem := range blueprints_data {
	// 		if blueprintItem.RepoName == state.RepoFilter.ValueString() {
	// 			filteredData = append(filteredData, blueprintItem)
	// 		}
	// 	}
	// } else {
	// 	filteredData = blueprints_data
	// }

	// for _, blueprintItem := range filteredData {
	// 	blueprintData := blueprintModel{
	// 		BlueprintName: types.StringValue(blueprintItem.BlueprintName),
	// 		Name:          types.StringValue(blueprintItem.Name),
	// 		RepoName:      types.StringValue(blueprintItem.RepoName),
	// 		Description:   types.StringValue(blueprintItem.Description),
	// 		Commit:        types.StringValue(blueprintItem.Commit),
	// 		ModifiedBy:    types.StringValue(blueprintItem.ModifiedBy),
	// 		DisplayName:   types.StringValue(blueprintItem.DisplayName),
	// 		RepoBranch:    types.StringValue(blueprintItem.RepoBranch),
	// 		Url:           types.StringValue(blueprintItem.Url),
	// 		Published:     types.BoolValue(blueprintItem.Published),
	// 	}
	// 	state.Blueprints = append(state.Blueprints, blueprintData)
	// }

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
