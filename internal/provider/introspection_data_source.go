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
	_ datasource.DataSource              = &environmentIntrospectionDataSource{}
	_ datasource.DataSourceWithConfigure = &environmentIntrospectionDataSource{}
)

// NewenvironmentIntrospectionsDataSource is a helper function to simplify the provider implementation.
func NewEnvironmentIntrospectionDataSource() datasource.DataSource {
	return &environmentIntrospectionDataSource{}
}

// environmentIntrospectionDataSource is the data source implementation.
type environmentIntrospectionDataSource struct {
	client *client.Client
}

// environmentIntrospectionDataSourceModel maps the data source schema data.
type environmentIntrospectionDataSourceModel struct {
	SpaceName        types.String        `tfsdk:"space_name"`
	Id               types.String        `tfsdk:"id"`
	GrainPath        types.String        `tfsdk:"grain_path"`
	GrainType        types.String        `tfsdk:"grain_type"`
	ResourceName     types.String        `tfsdk:"resource_name"`
	ResourceType     types.String        `tfsdk:"resource_type"`
	ResourceCategory types.String        `tfsdk:"resource_category"`
	Status           types.String        `tfsdk:"status"`
	Alias            types.String        `tfsdk:"alias"`
	HasRunningAction types.Bool          `tfsdk:"has_running_action"`
	Attributes       []keyValuePairModel `tfsdk:"attributes"`
	CustomIcon       types.String        `tfsdk:"custom_icon"`
}

// Metadata returns the data source type name.
func (d *environmentIntrospectionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment_introspection"
}

// Schema defines the schema for the data source.
func (d *environmentIntrospectionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get Environment Introspection details.",
		Attributes: map[string]schema.Attribute{
			"space_name": schema.StringAttribute{
				MarkdownDescription: "Torque's space this environmentIntrospection is in",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "EnvironmentIntrospection ID (15 alphanumeric characters)",
				Required:            true,
			},
			"grain_path": schema.StringAttribute{
				MarkdownDescription: "Name of the environmentIntrospection",
				Computed:            true,
			},
			"grain_type": schema.StringAttribute{
				MarkdownDescription: "Name of the blueprint that was used to launch this environmentIntrospection from",
				Computed:            true,
			},
			"resource_name": schema.StringAttribute{
				MarkdownDescription: "Short commit of the blueprint that was used to launch this environmentIntrospection from",
				Computed:            true,
			},
			"resource_type": schema.StringAttribute{
				MarkdownDescription: "Name of the blueprint's repository",
				Computed:            true,
			},
			"resource_category": schema.StringAttribute{
				MarkdownDescription: "EnvironmentIntrospection status",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Last time environmentIntrospection was accessed",
				Computed:            true,
			},
			"alias": schema.StringAttribute{
				MarkdownDescription: "Datetime string representing the time this nvironment was launched",
				Computed:            true,
			},
			"has_running_action": schema.BoolAttribute{
				MarkdownDescription: "Datetime string representing the time the environmentIntrospection has ended (if ended)",
				Computed:            true,
			},
			"custom_icon": schema.StringAttribute{
				MarkdownDescription: "Email address of the person who owns this environmentIntrospection",
				Computed:            true,
			},
			"attributes": schema.ListNestedAttribute{
				Description: "Actual inputs and their values that the environmentIntrospection was launched with",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Input's name",
							Computed:    true,
						},
						"value": schema.StringAttribute{
							Description: "The value of the input",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *environmentIntrospectionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *environmentIntrospectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state environmentIntrospectionDataSourceModel
	var space_name types.String
	var id types.String

	diags := req.Config.GetAttribute(ctx, path.Root("space_name"), &space_name)
	diags = append(diags, req.Config.GetAttribute(ctx, path.Root("id"), &id)...)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	introspection_data, err := d.client.GetIntrospectionDetails(space_name.ValueString(), id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Torque Environment Introspection",
			err.Error(),
		)
		return
	}
	for _, introspectionItem := range introspection_data {
		state.SpaceName = space_name
		state.Id = id
		state.GrainPath = types.StringValue(introspectionItem.GrainPath)
		state.GrainType = types.StringValue(introspectionItem.GrainType)
		state.ResourceName = types.StringValue(introspectionItem.ResourceName)
		state.ResourceType = types.StringValue(introspectionItem.ResourceType)
		state.ResourceCategory = types.StringValue(introspectionItem.ResourceCategory)
		state.Status = types.StringValue(introspectionItem.Status)
		state.Alias = types.StringValue(introspectionItem.Alias)
		state.HasRunningAction = types.BoolValue(introspectionItem.HasRunningAction)
		state.CustomIcon = types.StringValue(introspectionItem.CustomIcon)
		for _, attributeItem := range introspectionItem.Attributes {
			attributeData := keyValuePairModel{
				Name:  types.StringValue(attributeItem.Name),
				Value: types.StringValue(attributeItem.Value),
			}
			state.Attributes = append(state.Attributes, attributeData)
		}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
