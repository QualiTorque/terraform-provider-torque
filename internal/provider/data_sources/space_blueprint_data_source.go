package data_sources

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
	SpaceName               types.String     `tfsdk:"space_name"`
	Name                    types.String     `tfsdk:"name"`
	DisplayName             types.String     `tfsdk:"display_name"`
	RepoName                types.String     `tfsdk:"repository_name"`
	RepoBranch              types.String     `tfsdk:"repository_branch"`
	Commit                  types.String     `tfsdk:"commit"`
	Description             types.String     `tfsdk:"description"`
	Url                     types.String     `tfsdk:"url"`
	ModifiedBy              types.String     `tfsdk:"modified_by"`
	LastModified            types.String     `tfsdk:"last_modified"`
	Published               types.Bool       `tfsdk:"enabled"`
	NumOfActiveEnvironments types.Int32      `tfsdk:"num_of_active_environments"`
	Tags                    []blueprintTag   `tfsdk:"tags"`
	MaxDuration             types.String     `tfsdk:"max_duration"`
	DefaultDuration         types.String     `tfsdk:"default_duration"`
	DefaultExtend           types.String     `tfsdk:"default_extend"`
	MaxActiveEnvironments   types.Int32      `tfsdk:"max_active_environments"`
	AlwaysOn                types.Bool       `tfsdk:"always_on"`
	Inputs                  []blueprintInput `tfsdk:"inputs"`
}

type blueprintTag struct {
	Name           types.String `tfsdk:"name"`
	DefaultValue   types.String `tfsdk:"default_value"`
	PossibleValues types.List   `tfsdk:"possible_values"`
	Description    types.String `tfsdk:"description"`
}

type blueprintInput struct {
	Name           types.String `tfsdk:"name"`
	DefaultValue   types.String `tfsdk:"default_value"`
	PossibleValues types.List   `tfsdk:"possible_values"`
	Description    types.String `tfsdk:"description"`
}

// Metadata returns the data source type name.
func (d *spaceBlueprintDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_space_blueprint"
}

// Schema defines the schema for the data source.
func (d *spaceBlueprintDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Returns details of a blueprint in the specified space",
		Attributes: map[string]schema.Attribute{
			"space_name": schema.StringAttribute{
				MarkdownDescription: "Name of the space containing the blueprint",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Blueprint name",
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
				Description: "The time of the last modificiation of the blueprint",
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
							Description: "The tag's name",
							Computed:    true,
						},
						"default_value": schema.StringAttribute{
							Description: "The tag's default value",
							Computed:    true,
						},
						"possible_values": schema.ListAttribute{
							Description: "List of possible values for this tag",
							Computed:    true,
							ElementType: types.StringType,
						},
						"description": schema.StringAttribute{
							Description: "The description of this blueprint tag",
							Computed:    true,
						},
					},
				},
			},
			"inputs": schema.ListNestedAttribute{
				Description: "List of inputs that this blueprint requires.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Input's name",
							Computed:    true,
						},
						"default_value": schema.StringAttribute{
							Description: "Input's default value",
							Computed:    true,
						},
						"possible_values": schema.ListAttribute{
							Description: "List of possible values for this input",
							Computed:    true,
							ElementType: types.StringType,
						},
						"description": schema.StringAttribute{
							Description: "The input's description",
							Computed:    true,
						},
					},
				},
			},
			"num_of_active_environments": schema.Int32Attribute{
				Description: "Number of current active environments that were launched from this blueprint.",
				Computed:    true,
			},
			"max_duration": schema.StringAttribute{
				Description: "Max duration of environment in ISO 8601 format: P{days}DT{hours}H{minutes}M{seconds}S] (for example, 'P0DT2H3M4S')",
				Computed:    true,
			},
			"default_duration": schema.StringAttribute{
				Description: "Default duration of environment in ISO 8601 format: P{days}DT{hours}H{minutes}M{seconds}S] (for example, 'P0DT2H3M4S')",
				Computed:    true,
			},
			"default_extend": schema.StringAttribute{
				Description: "Default Extend of environment in ISO 8601 format: P{days}DT{hours}H{minutes}M{seconds}S] (for example, 'P0DT2H3M4S')",
				Computed:    true,
			},
			"max_active_environments": schema.Int32Attribute{
				Description: "Max active environments that can be launched from the blueprint",
				Computed:    true,
			},
			"always_on": schema.BoolAttribute{
				Description: "When True, the environment will always be on, when False, it will be ephemeral.",
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

	blueprint_data, err := d.client.GetBlueprint(space_name.ValueString(), name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Blueprint details or it doesn't exist",
			err.Error(),
		)
		return
	}

	// initialize state
	state.SpaceName = types.StringValue(space_name.ValueString())
	state.Name = types.StringValue(name.ValueString())
	state.DisplayName = types.StringValue(blueprint_data.DisplayName)
	state.Commit = types.StringValue(blueprint_data.Commit)
	state.Description = types.StringValue(blueprint_data.Description)
	state.Published = types.BoolValue(blueprint_data.Published)
	state.MaxDuration = types.StringValue(blueprint_data.Policies.MaxDuration)
	state.DefaultDuration = types.StringValue(blueprint_data.Policies.DefaultDuration)
	state.DefaultExtend = types.StringValue(blueprint_data.Policies.DefaultExtend)
	state.MaxActiveEnvironments = types.Int32Value(blueprint_data.Policies.MaxActiveEnvironments)
	state.AlwaysOn = types.BoolValue(blueprint_data.Policies.AlwaysOn)
	state.ModifiedBy = types.StringValue(blueprint_data.ModifiedBy)
	state.LastModified = types.StringValue(blueprint_data.LastModified)
	state.RepoBranch = types.StringValue(blueprint_data.RepoBranch)
	state.RepoName = types.StringValue(blueprint_data.RepoName)
	state.Url = types.StringValue(blueprint_data.Url)
	state.NumOfActiveEnvironments = types.Int32Value(blueprint_data.NumOfActiveEnvironments)
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

	for _, inputItem := range blueprint_data.Inputs {
		var possibleValues []attr.Value
		for _, value := range inputItem.PossibleValues {
			possibleValues = append(possibleValues, types.StringValue(value))
		}
		possibleValuesList, _ := types.ListValue(types.StringType, possibleValues)

		inputData := blueprintInput{
			Name:           types.StringValue(inputItem.Name),
			PossibleValues: possibleValuesList,
			DefaultValue:   types.StringValue(inputItem.DefaultValue),
			Description:    types.StringValue(inputItem.Description),
		}
		state.Inputs = append(state.Inputs, inputData)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
