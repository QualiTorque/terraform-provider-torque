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
	SpaceName               types.String        `tfsdk:"space_name"`
	Id                      types.String        `tfsdk:"id"`
	Name                    types.String        `tfsdk:"name"`
	IsEAC                   types.Bool          `tfsdk:"is_eac"`
	LastUsed                types.String        `tfsdk:"last_used"`
	BlueprintName           types.String        `tfsdk:"blueprint_name"`
	BlueprintCommit         types.String        `tfsdk:"blueprint_commit"`
	BlueprintRepositoryName types.String        `tfsdk:"blueprint_repository_name"`
	Status                  types.String        `tfsdk:"status"`
	OwnerEmail              types.String        `tfsdk:"owner_email"`
	InitiatorEmail          types.String        `tfsdk:"initiator_email"`
	StartTime               types.String        `tfsdk:"start_time"`
	EndTime                 types.String        `tfsdk:"end_time"`
	Errors                  []errorModel        `tfsdk:"errors"`
	Inputs                  []keyValuePairModel `tfsdk:"inputs"`
	Outputs                 []keyValuePairModel `tfsdk:"outputs"`
	Tags                    []keyValuePairModel `tfsdk:"tags"`
}

type keyValuePairModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type errorModel struct {
	Message types.String `tfsdk:"name"`
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

// type EnvironmentDetailTagsModel struct {
// 	Tags []TagModel `tfsdk:"tags"`
// }

// type TagModel struct {
// 	Name  types.String `tfsdk:"name"`
// 	Value types.String `tfsdk:"value"`
// }

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
			"name": schema.StringAttribute{
				MarkdownDescription: "Environment Name",
				Computed:            true,
			},
			"blueprint_name": schema.StringAttribute{
				MarkdownDescription: "Blueprint Name",
				Computed:            true,
			},
			"blueprint_commit": schema.StringAttribute{
				MarkdownDescription: "Blueprint Commit",
				Computed:            true,
			},
			"blueprint_repository_name": schema.StringAttribute{
				MarkdownDescription: "Name of the blueprint's repository",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Name of the blueprint's repository",
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
			"start_time": schema.StringAttribute{
				MarkdownDescription: "Last time environment was used",
				Computed:            true,
			},
			"end_time": schema.StringAttribute{
				MarkdownDescription: "Last time environment was used",
				Computed:            true,
			},
			"owner_email": schema.StringAttribute{
				MarkdownDescription: "Email address of the person who owns this environment",
				Computed:            true,
			},
			"initiator_email": schema.StringAttribute{
				MarkdownDescription: "Email address of the person who initiated this environment",
				Computed:            true,
			},
			"inputs": schema.ListNestedAttribute{
				Description: "Environment Inputs",
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
			"outputs": schema.ListNestedAttribute{
				Description: "Environment Inputs",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Output's name",
							Computed:    true,
						},
						"value": schema.StringAttribute{
							Description: "The value of the output",
							Computed:    true,
						},
					},
				},
			},
			"tags": schema.ListNestedAttribute{
				Description: "Environment Tags",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Tag's name",
							Computed:    true,
						},
						"value": schema.StringAttribute{
							Description: "The value of the tag",
							Computed:    true,
						},
					},
				},
			},
			"errors": schema.ListNestedAttribute{
				Description: "Environment Errors",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"message": schema.StringAttribute{
							Description: "Error Message",
							Computed:    true,
						},
					},
				},
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

	diags := req.Config.GetAttribute(ctx, path.Root("space_name"), &space_name)
	diags = append(diags, req.Config.GetAttribute(ctx, path.Root("id"), &id)...)
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

	state.Name = types.StringValue(environment_data.Details.Definition.Metadata.Name)
	state.LastUsed = types.StringValue(environment_data.LastUsed)
	state.IsEAC = types.BoolValue(environment_data.Details.State.IsEac)
	state.BlueprintName = types.StringValue(environment_data.Details.Definition.Metadata.BlueprintName)
	state.BlueprintCommit = types.StringValue(environment_data.Details.Definition.Metadata.BlueprintCommit)
	state.BlueprintRepositoryName = types.StringValue(environment_data.Details.Definition.Metadata.BlueprintRepositoryName)
	state.OwnerEmail = types.StringValue(environment_data.Owner.OwnerEmail)
	state.Status = types.StringValue(environment_data.Details.ComputedStatus)
	state.StartTime = types.StringValue(environment_data.Details.State.Execution.StartTime)
	state.EndTime = types.StringValue(environment_data.Details.State.Execution.EndTime)
	state.Id = types.StringValue(environment_data.Details.Id)
	state.InitiatorEmail = types.StringValue(environment_data.Initiator.InitiatorEmail)
	state.Inputs = []keyValuePairModel{}
	state.Tags = []keyValuePairModel{}
	state.Outputs = []keyValuePairModel{}
	state.Errors = []errorModel{}

	for _, inputItem := range environment_data.Details.Definition.Inputs {
		inputData := keyValuePairModel{
			Name:  types.StringValue(inputItem.Name),
			Value: types.StringValue(inputItem.Value),
		}
		state.Inputs = append(state.Inputs, inputData)
	}

	for _, tagItem := range environment_data.Details.Definition.Tags {
		tagData := keyValuePairModel{
			Name:  types.StringValue(tagItem.Name),
			Value: types.StringValue(tagItem.Value),
		}
		state.Tags = append(state.Tags, tagData)
	}
	for _, outputItem := range environment_data.Details.State.Outputs {
		outputData := keyValuePairModel{
			Name:  types.StringValue(outputItem.Name),
			Value: types.StringValue(outputItem.Value),
		}
		state.Outputs = append(state.Outputs, outputData)
	}

	for _, errorItem := range environment_data.Details.State.Errors {
		errorData := errorModel{
			Message: types.StringValue(errorItem.Message),
		}
		state.Errors = append(state.Errors, errorData)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
