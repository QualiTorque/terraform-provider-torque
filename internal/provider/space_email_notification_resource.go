package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/qualitorque/terraform-provider-torque/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TorqueSpaceEmailNotificationResource{}
var _ resource.ResourceWithImportState = &TorqueSpaceEmailNotificationResource{}

func NewTorqueSpaceEmailNotificationResource() resource.Resource {
	return &TorqueSpaceEmailNotificationResource{}
}

// TorqueSpaceEmailNotificationResource defines the resource implementation.
type TorqueSpaceEmailNotificationResource struct {
	client *client.Client
}

// TorqueSpaceEmailNotificationResourceModel describes the resource data model.
type TorqueSpaceEmailNotificationResourceModel struct {
	SpaceName                  types.String  `tfsdk:"space_name"`
	NotificationName           types.String  `tfsdk:"notification_name"`
	EnvironmentLaunched        types.Bool    `tfsdk:"environment_launched"`
	EnvironmentDeployed        types.Bool    `tfsdk:"environment_deployed"`
	EnvironmentForceEnded      types.Bool    `tfsdk:"environment_force_ended"`
	EnvironmentIdle            types.Bool    `tfsdk:"environment_idle"`
	EnvironmentExtended        types.Bool    `tfsdk:"environment_extended"`
	DriftDetected              types.Bool    `tfsdk:"drift_detected"`
	WorkflowFailed             types.Bool    `tfsdk:"workflow_failed"`
	WorkflowStarted            types.Bool    `tfsdk:"workflow_started"`
	UpdatesDetected            types.Bool    `tfsdk:"updates_detected"`
	CollaboratorAdded          types.Bool    `tfsdk:"collaborator_added"`
	ActionFailed               types.Bool    `tfsdk:"action_failed"`
	EnvironmentEndingFailed    types.Bool    `tfsdk:"environment_ending_failed"`
	EnvironmentEnded           types.Bool    `tfsdk:"environment_ended"`
	EnvironmentActiveWithError types.Bool    `tfsdk:"environment_active_with_error"`
	WorkflowStartReminder      types.Int64   `tfsdk:"workflow_start_reminder"` // if this is set - Kate - need to set also - "workflow_events_notifier": { "notify_on_all_workflows": true}
	EndThreashold              types.Int64   `tfsdk:"end_threshold"`
	IdleReminder               []types.Int64 `tfsdk:"idle_reminders"`
	NotificationId             types.String  `tfsdk:"notification_id"`
}

func (r *TorqueSpaceEmailNotificationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_space_email_notification"
}

func (r *TorqueSpaceEmailNotificationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Creation of a new notification is a Torque space",

		Attributes: map[string]schema.Attribute{
			"space_name": schema.StringAttribute{
				MarkdownDescription: "Space name to add the notification to",
				Required:            true,
			},
			"notification_name": schema.StringAttribute{
				MarkdownDescription: "The notification cofngiuration name in the space",
				Required:            true,
			},
			"environment_launched": schema.BoolAttribute{
				MarkdownDescription: "Configure notification for the \"Environment Launched\" event",
				Optional:            true,
				Computed:            false,
			},
			"environment_deployed": schema.BoolAttribute{
				MarkdownDescription: "Configure notification for the \"Environment Deployed\" event",
				Optional:            true,
				Computed:            false,
			},
			"environment_force_ended": schema.BoolAttribute{
				MarkdownDescription: "Configure notification for the \"Environment Force Ended\" event",
				Optional:            true,
				Computed:            false,
			},
			"environment_idle": schema.BoolAttribute{
				MarkdownDescription: "Configure notification for the \"Environment Idle\" event",
				Optional:            true,
				Computed:            false,
			},
			"environment_extended": schema.BoolAttribute{
				MarkdownDescription: "Configure notification for the \"Environment Extended\" event",
				Optional:            true,
				Computed:            false,
			},
			"drift_detected": schema.BoolAttribute{
				MarkdownDescription: "Configure notification for the \"Drift Detected\" event",
				Optional:            true,
				Computed:            false,
			},
			"workflow_failed": schema.BoolAttribute{
				MarkdownDescription: "Configure notification for the \"Workflow Failed\" event",
				Optional:            true,
				Computed:            false,
			},
			"workflow_started": schema.BoolAttribute{
				MarkdownDescription: "Configure notification for the \"Workflow Started\" event",
				Optional:            true,
				Computed:            false,
			},
			"updates_detected": schema.BoolAttribute{
				MarkdownDescription: "Configure notification for the \"Updates Detected\" event",
				Optional:            true,
				Computed:            false,
			},
			"collaborator_added": schema.BoolAttribute{
				MarkdownDescription: "Configure notification for the \"Collaborator Added\" event",
				Optional:            true,
				Computed:            false,
			},
			"action_failed": schema.BoolAttribute{
				MarkdownDescription: "Configure notification for the \"Action Failed\" event",
				Optional:            true,
				Computed:            false,
			},
			"environment_ending_failed": schema.BoolAttribute{
				MarkdownDescription: "Configure notification for the \"Environment Ending Failed\" event",
				Optional:            true,
				Computed:            false,
			},
			"environment_ended": schema.BoolAttribute{
				MarkdownDescription: "Configure notification for the \"Environment Ended\" event",
				Optional:            true,
				Computed:            false,
			},
			"environment_active_with_error": schema.BoolAttribute{
				MarkdownDescription: "Configure notification for the \"Environment Active With Error\" event",
				Optional:            true,
				Computed:            false,
			},
			"workflow_start_reminder": schema.Int64Attribute{
				MarkdownDescription: "Configure notification for the \"Drift Detected\" event",
				Optional:            true,
				Computed:            false,
			},
			"end_threshold": schema.Int64Attribute{
				MarkdownDescription: "Time in minutes to send notification environment end event reminder notification before an environment ends. For example, 10",
				Optional:            true,
				Computed:            false,
			},
			"idle_reminders": schema.ListAttribute{
				MarkdownDescription: "Array of time in hours to send notification for environment idle reminder",
				Optional:            true,
				Computed:            false,
				ElementType:         types.Int64Type,
			},
			"notification_id": schema.StringAttribute{
				MarkdownDescription: "The id of the newly added notification",
				Computed:            true,
			},
		},
	}
}

func (r *TorqueSpaceEmailNotificationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueSpaceEmailNotificationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueSpaceEmailNotificationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var idle []int64
	if len(data.IdleReminder) > 0 {
		for _, reminder := range data.IdleReminder {
			idle = append(idle, reminder.ValueInt64())
		}
	}

	notification, err := r.client.CreateSpaceEmailNotification(data.SpaceName.ValueString(), data.NotificationName.ValueString(), data.EnvironmentLaunched.ValueBool(),
		data.EnvironmentDeployed.ValueBool(), data.EnvironmentForceEnded.ValueBool(), data.EnvironmentIdle.ValueBool(), data.EnvironmentExtended.ValueBool(), data.DriftDetected.ValueBool(),
		data.WorkflowFailed.ValueBool(), data.WorkflowStarted.ValueBool(), data.UpdatesDetected.ValueBool(), data.CollaboratorAdded.ValueBool(), data.ActionFailed.ValueBool(),
		data.EnvironmentEndingFailed.ValueBool(), data.EnvironmentEnded.ValueBool(), data.EnvironmentActiveWithError.ValueBool(), data.WorkflowStartReminder.ValueInt64(), data.EndThreashold.ValueInt64(),
		idle)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create notification in space, got error: %s", err))
		return
	}

	data.NotificationId = basetypes.NewStringValue(strings.Replace(notification, "\"", "", -1))

	tflog.Trace(ctx, "Resource Created Successful!")

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueSpaceEmailNotificationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueSpaceEmailNotificationResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueSpaceEmailNotificationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueSpaceEmailNotificationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueSpaceEmailNotificationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueSpaceEmailNotificationResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the notification.
	err := r.client.DeleteSpaceNotification(data.SpaceName.ValueString(), data.NotificationId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete space notification, got error: %s", err))
		return
	}

}

func (r *TorqueSpaceEmailNotificationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("notification_name"), req, resp)
}
