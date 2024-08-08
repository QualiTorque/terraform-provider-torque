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
	CollaboratorsEmails     []collaboratorModel `tfsdk:"collaborators"`
	StartTime               types.String        `tfsdk:"start_time"`
	EndTime                 types.String        `tfsdk:"end_time"`
	Grains                  []grainModel        `tfsdk:"grains"`
	Errors                  []errorModel        `tfsdk:"errors"`
	Inputs                  []keyValuePairModel `tfsdk:"inputs"`
	Outputs                 []keyValuePairModel `tfsdk:"outputs"`
	Tags                    []keyValuePairModel `tfsdk:"tags"`
	RawJson                 types.String        `tfsdk:"raw_json"`
}

type keyValuePairModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type grainModel struct {
	Name    types.String       `tfsdk:"name"`
	Kind    types.String       `tfsdk:"kind"`
	Id      types.String       `tfsdk:"id"`
	Path    types.String       `tfsdk:"path"`
	State   grainStateModel    `tfsdk:"state"`
	Sources []grainSourceModel `tfsdk:"sources"`
}

type grainStateModel struct {
	CurrentState types.String `tfsdk:"current_state"`
}

type grainSourceModel struct {
	Store        types.String `tfsdk:"store"`
	Path         types.String `tfsdk:"path"`
	Branch       types.String `tfsdk:"branch"`
	Commit       types.String `tfsdk:"commit"`
	IsLastCommit types.Bool   `tfsdk:"is_last_commit"`
}

type errorModel struct {
	Message types.String `tfsdk:"name"`
}

type collaboratorModel struct {
	Email types.String `tfsdk:"email"`
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
				MarkdownDescription: "Torque's space this environment is in",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Environment ID (15 alphanumeric characters)",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the environment",
				Computed:            true,
			},
			"blueprint_name": schema.StringAttribute{
				MarkdownDescription: "Name of the blueprint that was used to launch this environment from",
				Computed:            true,
			},
			"blueprint_commit": schema.StringAttribute{
				MarkdownDescription: "Short commit of the blueprint that was used to launch this environment from",
				Computed:            true,
			},
			"blueprint_repository_name": schema.StringAttribute{
				MarkdownDescription: "Name of the blueprint's repository",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Environment status",
				Computed:            true,
			},
			"collaborators": schema.ListNestedAttribute{
				Description: "Environment collaborators",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"email": schema.StringAttribute{
							Description: "Collaborator's email",
							Computed:    true,
						},
					},
				},
			},
			"is_eac": schema.BoolAttribute{
				MarkdownDescription: "Is environment source is Env-as-Code",
				Computed:            true,
			},
			"last_used": schema.StringAttribute{
				MarkdownDescription: "Last time environment was accessed",
				Computed:            true,
			},
			"start_time": schema.StringAttribute{
				MarkdownDescription: "Datetime string representing the time this nvironment was launched",
				Computed:            true,
			},
			"end_time": schema.StringAttribute{
				MarkdownDescription: "Datetime string representing the time the environment has ended (if ended)",
				Computed:            true,
			},
			"owner_email": schema.StringAttribute{
				MarkdownDescription: "Email address of the person who owns this environment",
				Computed:            true,
			},
			"initiator_email": schema.StringAttribute{
				MarkdownDescription: "Email address of the person who initiated (launched) this environment",
				Computed:            true,
			},
			"raw_json": schema.StringAttribute{
				MarkdownDescription: "Raw JSON response",
				Computed:            true,
			},
			"grains": schema.ListNestedAttribute{
				Description: "Environment Inputs",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Grain's name",
							Computed:    true,
						},
						"kind": schema.StringAttribute{
							Description: "Grain's Kind",
							Computed:    true,
						},
						"id": schema.StringAttribute{
							Description: "Grain's id",
							Computed:    true,
						},
						"path": schema.StringAttribute{
							Description: "Grain's path in the repository (store)",
							Computed:    true,
						},
						"state": schema.SingleNestedAttribute{
							MarkdownDescription: "Additional details about the environment state.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"current_state": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "Grain's state",
								},
							},
						},
						"sources": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"store": schema.StringAttribute{
										MarkdownDescription: "The store (repository) of this grain",
										Computed:            true,
									},
									"path": schema.StringAttribute{
										MarkdownDescription: "The path in the repository (store)",
										Computed:            true,
									},
									"branch": schema.StringAttribute{
										MarkdownDescription: "The branch used as the source",
										Computed:            true,
									},
									"commit": schema.StringAttribute{
										MarkdownDescription: "The commit used as the sorce",
										Computed:            true,
									},
									"is_last_commit": schema.BoolAttribute{
										MarkdownDescription: "Specify whether the commit is the latest of the source",
										Computed:            true,
									},
								},
							},
						},
					},
				},
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

	environment_data, raw_json, err := d.client.GetEnvironmentDetails(space_name.ValueString(), id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Torque environment",
			err.Error(),
		)
		return
	}

	state.Name = types.StringValue(environment_data.Details.Definition.Metadata.Name)
	state.SpaceName = types.StringValue(environment_data.Details.Definition.Metadata.SpaceName)
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
	state.RawJson = types.StringValue(raw_json)

	state.CollaboratorsEmails = []collaboratorModel{}
	state.Inputs = []keyValuePairModel{}
	state.Tags = []keyValuePairModel{}
	state.Outputs = []keyValuePairModel{}
	state.Errors = []errorModel{}

	for _, grainItem := range environment_data.Details.State.Grains {
		grainData := grainModel{
			Name: types.StringValue(grainItem.Name),
			Kind: types.StringValue(grainItem.Kind),
			Id:   types.StringValue(grainItem.Id),
			Path: types.StringValue(grainItem.Path),
			State: grainStateModel{
				CurrentState: types.StringValue(grainItem.State.CurrentState),
			},
		}
		for _, grainSource := range grainItem.Sources {
			grainSourceData := grainSourceModel{
				Store:        types.StringValue(grainSource.Store),
				Path:         types.StringValue(grainSource.Path),
				Branch:       types.StringValue(grainSource.Branch),
				Commit:       types.StringValue(grainSource.Commit),
				IsLastCommit: types.BoolValue(grainSource.IsLastCommit),
			}
			grainData.Sources = append(grainData.Sources, grainSourceData)
		}

		state.Grains = append(state.Grains, grainData)
	}
	for _, collaboratorItem := range environment_data.CollaboratorsInfo.Collaborators {
		collaboratorData := collaboratorModel{
			Email: types.StringValue(collaboratorItem.Email),
		}
		state.CollaboratorsEmails = append(state.CollaboratorsEmails, collaboratorData)
	}

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
