package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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

type TorqueEnvironmentResourceModel struct {
	EnvironmentName types.String `tfsdk:"environment_name"`
	BlueprintName   types.String `tfsdk:"blueprint_name"`
	Space           types.String `tfsdk:"space"`
	Id              types.String `tfsdk:"id"`
	// OwnerEmail       types.String         `tfsdk:"owner_email"`
	// Description      types.String         `tfsdk:"description"`
	// Inputs           types.Map            `tfsdk:"inputs"`
	// Tags             types.Map            `tfsdk:"tags"`
	// Collaborators    []CollaboratorsModel `tfsdk:"collaborators"`
	// Automation       types.Bool           `tfsdk:"automation"`
	// ScheduledEndTime types.String         `tfsdk:"scheduled_end_time"`
	Duration types.String `tfsdk:"duration"`
	// Source           struct {
	// 	BlueprintName  types.String `tfsdk:"blueprint_name"`
	// 	RepositoryName types.String `tfsdk:"repository_name"`
	// 	Branch         types.String `tfsdk:"branch"`
	// 	Commit         types.String `tfsdk:"commit"`
	// } `tfsdk:"source"`
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
				Required:            true,
				Computed:            false,
			},
			// "inputs": schema.MapAttribute{
			// 	MarkdownDescription: "A list of inputs",
			// 	ElementType:         types.StringType,
			// 	Required:            false,
			// 	Computed:            false,
			// 	Optional:            true,
			// },
			// "description": schema.StringAttribute{
			// 	MarkdownDescription: "A list of inputs",
			// 	Required:            false,
			// 	Computed:            false,
			// 	Optional:            true,
			// },
			// "tags": schema.MapAttribute{
			// 	MarkdownDescription: "A list of inputs",
			// 	ElementType:         types.StringType,
			// 	Required:            false,
			// 	Computed:            true,
			// 	Optional:            true,
			// },
			// "collaborators": schema.ListNestedAttribute{
			// 	Description: "key-value pairs of spaces and roles that the newly created group will be associated to",
			// 	Computed:    false,
			// 	Optional:    true,
			// 	NestedObject: schema.NestedAttributeObject{
			// 		Attributes: map[string]schema.Attribute{
			// 			"collaborators_emails": schema.StringAttribute{
			// 				Description: "An existing Torque space name",
			// 				Computed:    false,
			// 				Optional:    true,
			// 			},
			// 			"all_space_members": schema.BoolAttribute{
			// 				Description: "Space role to be used for the specific space in the group",
			// 				Computed:    false,
			// 				Optional:    true,
			// 			},
			// 		},
			// 	},
			// },

			// "source": schema.StringAttribute{
			// 	MarkdownDescription: "A list of inputs",
			// 	Required:            false,
			// 	Computed:            false,
			// 	Optional:            true,
			// },
			"space": schema.StringAttribute{
				MarkdownDescription: "A list of inputs",
				Required:            false,
				Computed:            false,
				Optional:            true,
			},
			// "scheduled_end_time": schema.StringAttribute{
			// 	MarkdownDescription: "A list of inputs",
			// 	Required:            false,
			// 	Computed:            true,
			// 	Optional:            true,
			// },
			// "automation": schema.BoolAttribute{
			// 	MarkdownDescription: "A list of inputs",
			// 	Required:            false,
			// 	Computed:            false,
			// 	Optional:            true,
			// },
			"id": schema.StringAttribute{
				MarkdownDescription: "A list of inputs",
				Required:            false,
				Computed:            true,
			},
			// "owner_email": schema.StringAttribute{
			// 	MarkdownDescription: "A list of inputs",
			// 	Required:            false,
			// 	Computed:            true,
			// },
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

	body, err := r.client.CreateEnvironment(data.Space.ValueString(), data.BlueprintName.ValueString(), data.EnvironmentName.ValueString(), data.Duration.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Environment, got error: %s", err))
		return
	}
	// Initialize the inputs map
	// var inputs = make(map[string]string)

	// if !data.Inputs.IsNull() {
	// 	for key, value := range data.Inputs.Elements() {
	// 		inputs[key] = value.String()
	// 	}
	// }
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
