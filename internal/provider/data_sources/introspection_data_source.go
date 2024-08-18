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

type introspectionDataModel struct {
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

// environmentIntrospectionDataSourceModel maps the data source schema data.
type environmentIntrospectionDataSourceModel struct {
	SpaceName types.String             `tfsdk:"space_name"`
	Id        types.String             `tfsdk:"id"`
	Resources []introspectionDataModel `tfsdk:"resources"`
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
				MarkdownDescription: "Torque's space this environment is in",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Environment ID (15 alphanumeric characters)",
				Required:            true,
			},
			"resources": schema.ListNestedAttribute{
				Description: "List of introspection data, the actual resources in the environment",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"grain_path": schema.StringAttribute{
							MarkdownDescription: "Grain's path in the blueprint",
							Computed:            true,
						},
						"grain_type": schema.StringAttribute{
							MarkdownDescription: "Grain kind, like Terraform, CloudFormation, etc.",
							Computed:            true,
						},
						"resource_name": schema.StringAttribute{
							MarkdownDescription: "Resource name as it's named in the Infra-as-Code assets (Terraform, OpenTofu, etc.)",
							Computed:            true,
						},
						"resource_type": schema.StringAttribute{
							MarkdownDescription: "Resource type, like azure_vm or_aws instance for example",
							Computed:            true,
						},
						"resource_category": schema.StringAttribute{
							MarkdownDescription: "Category of resource, for example, storage, compute etc.",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							MarkdownDescription: "Status of the resource, for example ec2 instance might be running or shut down etc.",
							Computed:            true,
						},
						"alias": schema.StringAttribute{
							MarkdownDescription: "Resource concrete name",
							Computed:            true,
						},
						"has_running_action": schema.BoolAttribute{
							MarkdownDescription: "Determines if there's an action/workflow running on this resource",
							Computed:            true,
						},
						"custom_icon": schema.StringAttribute{
							MarkdownDescription: "Path to a custom icon for the resource",
							Computed:            true,
						},
						"attributes": schema.ListNestedAttribute{
							Description: "Map of the resource attributes",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "Attribute's name",
										Computed:    true,
									},
									"value": schema.StringAttribute{
										Description: "The value of the Attribute",
										Computed:    true,
									},
								},
							},
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
	state.SpaceName = space_name
	state.Id = id
	for _, introspectionItem := range introspection_data {
		introspectionData := introspectionDataModel{
			GrainPath:        types.StringValue(introspectionItem.GrainPath),
			GrainType:        types.StringValue(introspectionItem.GrainType),
			ResourceName:     types.StringValue(introspectionItem.ResourceName),
			ResourceType:     types.StringValue(introspectionItem.ResourceType),
			ResourceCategory: types.StringValue(introspectionItem.ResourceCategory),
			Status:           types.StringValue(introspectionItem.Status),
			Alias:            types.StringValue(introspectionItem.Alias),
			HasRunningAction: types.BoolValue(introspectionItem.HasRunningAction),
			CustomIcon:       types.StringValue(introspectionItem.CustomIcon),
		}
		for _, attributeItem := range introspectionItem.Attributes {
			attributeData := keyValuePairModel{
				Name:  types.StringValue(attributeItem.Name),
				Value: types.StringValue(attributeItem.Value),
			}
			introspectionData.Attributes = append(introspectionData.Attributes, attributeData)
		}
		state.Resources = append(state.Resources, introspectionData)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
