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

type SourceModel struct {
	BlueprintName  *string `tfsdk:"blueprint_name"`
	RepositoryName *string `tfsdk:"repository_name"`
	Branch         *string `tfsdk:"branch"`
	Commit         *string `tfsdk:"commit"`
}

type TorqueEnvironmentResourceModel struct {
	EnvironmentName  types.String        `tfsdk:"environment_name"`
	BlueprintName    types.String        `tfsdk:"blueprint_name"`
	Space            types.String        `tfsdk:"space"`
	Id               types.String        `tfsdk:"id"`
	OwnerEmail       types.String        `tfsdk:"owner_email"`
	Description      types.String        `tfsdk:"description"`
	Inputs           types.Map           `tfsdk:"inputs"`
	Tags             types.Map           `tfsdk:"tags"`
	Collaborators    *CollaboratorsModel `tfsdk:"collaborators"`
	Automation       types.Bool          `tfsdk:"automation"`
	ScheduledEndTime types.String        `tfsdk:"scheduled_end_time"`
	Duration         types.String        `tfsdk:"duration"`
	Source           *SourceModel        `tfsdk:"source"`
	// Workflows []struct {
	// 	Name      string `tfsdk:"name"`
	// 	Schedules []struct {
	// 		Scheduler  string     `tfsdk:"scheduler"`
	// 		Overridden types.Bool `tfsdk:"overridden"`
	// 	} `tfsdk:"schedules"`
	// 	Reminder        types.String      `tfsdk:"reminder"`
	// 	InputsOverrides map[string]string `tfsdk:"inputs_overrides"`
	// } `tfsdk:"workflows"`
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
				MarkdownDescription: "Name of the new Torque environment.",
				Required:            true,
			},
			"duration": schema.StringAttribute{
				MarkdownDescription: "Duration of environment",
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
				MarkdownDescription: "A list of inputs",
				ElementType:         types.StringType,
				Required:            false,
				Computed:            false,
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A list of inputs",
				Required:            false,
				Computed:            false,
				Optional:            true,
			},
			"tags": schema.MapAttribute{
				MarkdownDescription: "Environment blueprint tags",
				ElementType:         types.StringType,
				Required:            false,
				Computed:            false,
				Optional:            true,
			},
			"collaborators": schema.ObjectAttribute{
				Description: "key-value pairs of spaces and roles that the newly created group will be associated to",
				Computed:    false,
				Optional:    true,
				Required:    false,
				AttributeTypes: map[string]attr.Type{
					"collaborators_emails": types.ListType{
						ElemType: types.StringType,
					},
					"all_space_members": types.BoolType,
				},
			},
			"source": schema.SingleNestedAttribute{
				MarkdownDescription: "Additional details about the blueprint repository to be used. By default, this information is taken from the repository already confiured in the space.",
				Required:            false,
				Computed:            false,
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"blueprint_name": schema.StringAttribute{
						Optional: true,
					},
					"repository_name": schema.StringAttribute{
						Optional: true,
					},
					"branch": schema.StringAttribute{
						Optional: true,
					},
					"commit": schema.StringAttribute{
						Optional: true,
					},
				},
				// AttributeTypes: map[string]schema.Attribute{
				// 	"blueprint_name": schema.StringAttribute{
				// 		Description: "An existing Torque space name",
				// 		Computed:    false,
				// 		Optional:    true,
				// 	},
				// 	"repository_name": schema.StringAttribute{
				// 		Description: "Space role to be used for the specific space in the group",
				// 		Computed:    false,
				// 		Optional:    true,
				// 	},
				// 	"branch": schema.StringAttribute{
				// 		Description: "An existing Torque space name",
				// 		Computed:    false,
				// 		Optional:    true,
				// 	},
				// 	"commit": schema.StringAttribute{
				// 		Description: "Space role to be used for the specific space in the group",
				// 		Computed:    false,
				// 		Optional:    true,
				// 	},
				// },
			},
			"space": schema.StringAttribute{
				MarkdownDescription: "A list of inputs",
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
				MarkdownDescription: "A list of inputs",
				Required:            false,
				Computed:            true,
			},
			"owner_email": schema.StringAttribute{
				MarkdownDescription: "A list of inputs",
				Required:            false,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("someemail@quali.com"),
			},
			// "workflows": schema.StringAttribute{
			// 	MarkdownDescription: "A list of inputs",
			// 	Required:            false,
			// 	Computed:            false,
			// 	Optional:            true,
			// },
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
			inputs[key] = value.String()
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
	var source client.Source
	if data.Source != nil {
		if data.Source.BlueprintName != nil {
			source.BlueprintName = data.Source.BlueprintName
		}
		if data.Source.RepositoryName != nil {
			source.RepositoryName = data.Source.RepositoryName
		}
		if data.Source.Branch != nil {
			source.Branch = data.Source.Branch
		}
		if data.Source.Commit != nil {
			source.Commit = data.Source.Commit
		}
	}
	body, err := r.client.CreateEnvironment(data.Space.ValueString(), data.BlueprintName.ValueString(), data.EnvironmentName.ValueString(), data.Duration.ValueString(),
		inputs, data.OwnerEmail.ValueString(), data.Automation.ValueBool(), tags, collaborators, data.ScheduledEndTime.ValueString(), source)
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
