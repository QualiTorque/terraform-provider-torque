package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/qualitorque/terraform-provider-torque/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TorqueEnvironmentResource{}
var _ resource.ResourceWithImportState = &TorqueEnvironmentResource{}

func NewTorqueEnvironmentResource() resource.Resource {
	return &TorqueEnvironmentResource{}
}

// TorqueEnvironmentResource defines the resource implementation.
type TorqueEnvironmentResource struct {
	client *client.Client
}

type CollaboratorsModel struct {
	CollaboratorsEmails types.List `tfsdk:"collaborators_emails"`
	AllSpaceMembers     types.Bool `tfsdk:"all_space_members"`
}

type BlueprintSourceModel struct {
	BlueprintName  *string `tfsdk:"blueprint_name"`
	RepositoryName *string `tfsdk:"repository_name"`
	Branch         *string `tfsdk:"branch"`
	Commit         *string `tfsdk:"commit"`
}

type WorkflowModel struct {
	Name            types.String    `tfsdk:"name"`
	Schedules       []ScheduleModel `tfsdk:"schedules"`
	Reminder        types.Int64     `tfsdk:"reminder"`
	InputsOverrides types.Map       `tfsdk:"inputs_overrides"`
}
type ScheduleModel struct {
	Scheduler  types.String `tfsdk:"scheduler"`
	Overridden types.Bool   `tfsdk:"overridden"`
}

type TorqueEnvironmentResourceModel struct {
	EnvironmentName  types.String          `tfsdk:"environment_name"`
	BlueprintName    types.String          `tfsdk:"blueprint_name"`
	Space            types.String          `tfsdk:"space"`
	Id               types.String          `tfsdk:"id"`
	OwnerEmail       types.String          `tfsdk:"owner_email"`
	Description      types.String          `tfsdk:"description"`
	Inputs           types.Map             `tfsdk:"inputs"`
	Tags             types.Map             `tfsdk:"tags"`
	Collaborators    *CollaboratorsModel   `tfsdk:"collaborators"`
	Automation       types.Bool            `tfsdk:"automation"`
	ScheduledEndTime types.String          `tfsdk:"scheduled_end_time"`
	Duration         types.String          `tfsdk:"duration"`
	BlueprintSource  *BlueprintSourceModel `tfsdk:"blueprint_source"`
	Workflows        []WorkflowModel       `tfsdk:"workflows"`
}

func (r *TorqueEnvironmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_environment"
}

func (r *TorqueEnvironmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Creation of a new Torque Environment",

		Attributes: map[string]schema.Attribute{
			"blueprint_name": schema.StringAttribute{
				MarkdownDescription: "Name of the Torque blueprint that the torque environment will be launched from.",
				Required:            true,
			},
			"environment_name": schema.StringAttribute{
				MarkdownDescription: "The name for the newly created environment. Environment name can contain any character including special character and spaces.",
				Required:            true,
			},
			"duration": schema.StringAttribute{
				MarkdownDescription: "Environment duration time in ISO 8601 format: 'P{days}DT{hours}H{minutes}M{seconds}S]]' For example, P0DT2H3M4S. NOTE: Environment request cannot include both 'duration' and 'scheduled_end_time' fields.",
				Optional:            true,
				Computed:            false,
				Validators: []validator.String{
					// Validate only this attribute or other_attr is configured.
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("scheduled_end_time"),
					}...),
				},
			},
			"inputs": schema.MapAttribute{
				MarkdownDescription: "Dictionary of key-value string pairs that will be used as values for the blueprint inputs. In case a value is not provided the input default value will be used. If a default value is not set, a validation error will be thrown upon launch. For example: { 'region': 'eu-west-1', 'application version': '1.0.8' }",
				ElementType:         types.StringType,
				Required:            false,
				Computed:            false,
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The new environment description that will be presented in the Torque UI following the launch of the environment.",
				Required:            false,
				Computed:            false,
				Optional:            true,
			},
			"tags": schema.MapAttribute{
				MarkdownDescription: "Environment blueprint tags /// Dictionary of key-value string pairs that will be used to tag deployed resources in the environment. In case a configured tag value is not provided the tag default value will be used. Note that tags that were configured in the account and space level will be set regardless of this field. For example: { 'activity_type': 'demo'}",
				ElementType:         types.StringType,
				Required:            false,
				Computed:            false,
				Optional:            true,
			},
			"collaborators": schema.ObjectAttribute{
				MarkdownDescription: "Object of collaborators to add to the environment. Provide collaborators_emails list of strings representing emails of users in the account or set all_space_users to true to add everyone in the space",
				Computed:            false,
				Optional:            true,
				Required:            false,
				AttributeTypes: map[string]attr.Type{
					"collaborators_emails": types.ListType{
						ElemType: types.StringType,
					},
					"all_space_members": types.BoolType,
				},
			},
			"blueprint_source": schema.SingleNestedAttribute{
				MarkdownDescription: "Additional details about the blueprint repository to be used. By default, this information is taken from the repository already confiured in the space.",
				Required:            false,
				Computed:            false,
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"blueprint_name": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Sandbox blueprint name",
					},
					"repository_name": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The name of the repo to be used. This repository should be on-boarded to the space prior to launching the environment. In case you want to launch a 'Stored in Torque' Blueprint, you should set this field to 'qtorque'",
					},
					"branch": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Use this field to provide a branch from which the blueprint yaml will be launched",
					},
					"commit": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Use this field to provide a specific commit id from which the blueprint yaml will be launched",
					},
				},
			},
			"space": schema.StringAttribute{
				MarkdownDescription: "The space where this environment will be launched",
				Required:            false,
				Computed:            false,
				Optional:            true,
			},
			"scheduled_end_time": schema.StringAttribute{
				MarkdownDescription: "Environment scheduled end time in ISO 8601 format For example, 2021-10-06T08:27:05.215Z. NOTE: Environment request cannot include both 'duration' and 'scheduled_end_time' fields.",
				Computed:            false,
				Optional:            true,
				Validators: []validator.String{
					// Validate only this attribute or other_attr is configured.
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("duration"),
					}...),
				},
			},
			"automation": schema.BoolAttribute{
				MarkdownDescription: "Indicates if the environment was launched from automation using integrated pipeline tool, For example: Jenkins, GitHub Actions and GitLal CI.",
				Required:            false,
				Computed:            false,
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the environment",
				Required:            false,
				Computed:            true,
			},
			"owner_email": schema.StringAttribute{
				MarkdownDescription: "The email of the user that should be set as the owner of the new environment. if omitted the current user will be used.",
				Required:            false,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("someemail@quali.com"),
			},
			"workflows": schema.ListNestedAttribute{
				MarkdownDescription: "Array of workflows that will be attached and enabled on the new environment.",
				Required:            false,
				Computed:            false,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of existing and enabled workflow in the space",
							Computed:    false,
							Optional:    true,
						},
						"schedules": schema.ListNestedAttribute{
							Required: false,
							Computed: false,
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"scheduler": schema.StringAttribute{
										MarkdownDescription: "The CRON expression that schedules this workflow",
										Computed:            false,
										Optional:            true,
									},
									"overridden": schema.BoolAttribute{
										MarkdownDescription: "Specify if the workflow schedule can be overridden at launch",
										Computed:            false,
										Optional:            true,
									},
								},
							},
						},
						"reminder": schema.Int64Attribute{
							MarkdownDescription: "",
							Computed:            false,
							Optional:            true,
						},
						"inputs_overrides": schema.MapAttribute{
							MarkdownDescription: "Dictionary of key-value string pairs that can override the blueprint inputs ",
							ElementType:         types.StringType,
							Required:            false,
							Computed:            false,
							Optional:            true,
						},
					},
				},
			},
		},
	}
}

func (r *TorqueEnvironmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *TorqueEnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueEnvironmentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Initialize the inputs map
	var inputs = make(map[string]string)

	if !data.Inputs.IsNull() {
		for key, value := range data.Inputs.Elements() {
			inputs[key] = strings.Replace(value.String(), "\"", "", -1)
		}
	}
	var tags = make(map[string]string)

	if !data.Tags.IsNull() {
		for key, value := range data.Tags.Elements() {
			tags[key] = strings.Replace(value.String(), "\"", "", -1)
		}
	}
	var collaborators client.Collaborators
	if data.Collaborators != nil {
		var emails []string
		for _, val := range data.Collaborators.CollaboratorsEmails.Elements() {
			emails = append(emails, strings.Replace(val.String(), "\"", "", -1))
		}
		collaborators.AllSpaceMembers = data.Collaborators.AllSpaceMembers.ValueBool()
		collaborators.CollaboratorsEmails = emails
	}
	var blueprint_source client.BlueprintSource
	if data.BlueprintSource != nil {
		if data.BlueprintSource.BlueprintName != nil {
			blueprint_source.BlueprintName = data.BlueprintSource.BlueprintName
		}
		if data.BlueprintSource.RepositoryName != nil {
			blueprint_source.RepositoryName = data.BlueprintSource.RepositoryName
		}
		if data.BlueprintSource.Branch != nil {
			blueprint_source.Branch = data.BlueprintSource.Branch
		}
		if data.BlueprintSource.Commit != nil {
			blueprint_source.Commit = data.BlueprintSource.Commit
		}
	}
	var workflows []client.Workflow
	var inputs_overrides = make(map[string]string)
	var schedules []client.Schedule
	if len(data.Workflows) > 0 {
		for _, workflow := range data.Workflows {
			if len(workflow.Schedules) > 0 {
				for _, schedule := range workflow.Schedules {
					schedules = append(schedules, client.Schedule{
						Scheduler:  schedule.Scheduler.ValueString(),
						Overridden: schedule.Overridden.ValueBool(),
					})
				}
			}
			workflows = append(workflows, client.Workflow{
				Name:            workflow.Name.ValueString(),
				Reminder:        workflow.Reminder.ValueInt64(),
				InputsOverrides: inputs_overrides,
				Schedules:       schedules,
			})
		}
	}

	body, err := r.client.CreateEnvironment(data.Space.ValueString(), data.BlueprintName.ValueString(), data.EnvironmentName.ValueString(), data.Duration.ValueString(), data.Description.ValueString(),
		inputs, data.OwnerEmail.ValueString(), data.Automation.ValueBool(), tags, collaborators, data.ScheduledEndTime.ValueString(), blueprint_source, workflows)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Environment, got error: %s", err))
		return
	}

	var responseBody map[string]string
	if err := json.Unmarshal(body, &responseBody); err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Failed to parse response body: %s", err))
		return
	}

	id, ok := responseBody["id"]
	if !ok {
		resp.Diagnostics.AddError("ID Error", "ID not found in response body or is not of type string")
		return
	}
	data.Id = types.StringValue(id)

	// owner_email, ok := responseBody["owner_email"]
	// if !ok {
	// 	resp.Diagnostics.AddError("Owner email error", "Owner does not exist")
	// 	return
	// }
	// if owner_email != "" {
	// 	data.OwnerEmail = types.StringValue(owner_email)
	// }

	tflog.Trace(ctx, "Resource Created Successful!")

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueEnvironmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueEnvironmentResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueEnvironmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// var data TorqueEnvironmentResourceModel

	// Read Terraform plan data into the model
	// resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// if resp.Diagnostics.HasError() {
	// 	return
	// }
	// // If applicable, this is a great opportunity to initialize any necessary
	// // provider client data and make a call using it.
	// // httpResp, err := r.client.Do(httpReq)
	// // if err != nil {
	// //     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	// //     return
	// // }

	// // Save updated data into Terraform state
	// resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	resp.Diagnostics.AddError(
		"Resource updates of resource type torque_account are not permitted",
		"Cannot change details of torque account, use terraform destroy to delete it or create a new one",
	)
}

func (r *TorqueEnvironmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueEnvironmentResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Terminate the Environment.
	err := r.client.TerminateEnvironment(data.Space.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Environment, got error: %s", err))
		return
	}

}

func (r *TorqueEnvironmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
