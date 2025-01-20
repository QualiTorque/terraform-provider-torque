package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/qualitorque/terraform-provider-torque/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TorqueEmailApprovalChannelResource{}
var _ resource.ResourceWithImportState = &TorqueEmailApprovalChannelResource{}

func NewTorqueEmailApprovalChannelResource() resource.Resource {
	return &TorqueEmailApprovalChannelResource{}
}

// TorqueEmailApprovalChannelResource defines the resource implementation.
type TorqueEmailApprovalChannelResource struct {
	client *client.Client
}

// TorqueEmailApprovalChannelResourceModel describes the resource data model.
type TorqueEmailApprovalChannelResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Approvers   types.List   `tfsdk:"approvers"`
}

const (
	approval_channel_type = "Email"
)

func (r *TorqueEmailApprovalChannelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_email_approval_channel"
}

func (r *TorqueEmailApprovalChannelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Creates a new email approval channel.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the approval channel.",
				Optional:            false,
				Computed:            false,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the approval channel",
				Optional:            true,
				Computed:            true,
				Required:            false,
				Default:             stringdefault.StaticString(""),
			},
			"approvers": schema.ListAttribute{
				Description: "List of existing emails of users that will be the approvers of this approval channel",
				Required:    true,
				Optional:    false,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1), // Ensure the list has at least one entry if required
				},
			},
		},
	}
}
func (r *TorqueEmailApprovalChannelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueEmailApprovalChannelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueEmailApprovalChannelResourceModel
	var details client.ApprovalChannelDetails
	var approvers []client.Approver
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	for _, approver := range data.Approvers.Elements() {
		approvers = append(approvers, client.Approver{
			UserEmail: strings.Replace(approver.String(), "\"", "", -1),
		})
	}
	details.Approvers = approvers
	details.Type = approval_channel_type
	err := r.client.CreateApprovalChannel(data.Name.ValueString(), data.Description.ValueString(), details)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Approval Channel, got error: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueEmailApprovalChannelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueEmailApprovalChannelResourceModel
	approvers := []string{}
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	approval_channel, err := r.client.GetApprovalChannel(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Approval Channel details",
			"Could not read Approval Channel "+data.Name.ValueString()+": "+err.Error(),
		)
		return
	}
	data.Name = types.StringValue(approval_channel.Name)
	data.Description = types.StringValue(approval_channel.Description)
	for _, approver := range approval_channel.Details.Approvers {
		approvers = append(approvers, approver.UserEmail)
	}
	data.Approvers, _ = types.ListValueFrom(ctx, types.StringType, approvers)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TorqueEmailApprovalChannelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueEmailApprovalChannelResourceModel
	var details client.ApprovalChannelDetails
	var approvers []client.Approver

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	for _, approver := range data.Approvers.Elements() {
		approvers = append(approvers, client.Approver{
			UserEmail: strings.Replace(approver.String(), "\"", "", -1),
		})
	}
	details.Approvers = approvers
	details.Type = approval_channel_type
	err := r.client.UpdateApprovalChannel(data.Name.ValueString(), data.Description.ValueString(), details)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update Input Source, got error: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueEmailApprovalChannelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueEmailApprovalChannelResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteApprovalChannel(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Approval Channel, got error: %s", err))
		return
	}

}

func (r *TorqueEmailApprovalChannelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
