package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/qualitorque/terraform-provider-torque/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TorqueServiceNowApprovalChannelResource{}
var _ resource.ResourceWithImportState = &TorqueServiceNowApprovalChannelResource{}

func NewTorqueServiceNowApprovalChannelResource() resource.Resource {
	return &TorqueServiceNowApprovalChannelResource{}
}

// TorqueServiceNowApprovalChannelResource defines the resource implementation.
type TorqueServiceNowApprovalChannelResource struct {
	client *client.Client
}

// TorqueServiceNowApprovalChannelResourceModel describes the resource data model.
type TorqueServiceNowApprovalChannelResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Approver    types.String `tfsdk:"approver"`
	BaseUrl     types.String `tfsdk:"base_url"`
	UserName    types.String `tfsdk:"user_name"`
	Password    types.String `tfsdk:"password"`
	Headers     types.String `tfsdk:"headers"`
}

const (
	servicenow_approval_channel_type = "ServiceNow"
)

func (r *TorqueServiceNowApprovalChannelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_servicenow_approval_channel"
}

func (r *TorqueServiceNowApprovalChannelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Creates a new ServiceNow approval channel.",
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
			"base_url": schema.StringAttribute{
				MarkdownDescription: "ServiceNow Instance Base URL",
				Optional:            false,
				Computed:            false,
				Required:            true,
			},
			"user_name": schema.StringAttribute{
				MarkdownDescription: "ServiceNow Username",
				Optional:            false,
				Computed:            false,
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "ServiceNow Password",
				Optional:            false,
				Computed:            false,
				Required:            true,
				Sensitive:           true,
			},
			"headers": schema.StringAttribute{
				MarkdownDescription: "Custom Headers (JSON) - JSON formatted string that represents the custom headers, for example {header:'val'}",
				Optional:            true,
				Computed:            false,
				Required:            false,
			},
			"approver": schema.StringAttribute{
				MarkdownDescription: "ServiceNow Approver",
				Optional:            false,
				Computed:            false,
				Required:            true,
			},
		},
	}
}
func (r *TorqueServiceNowApprovalChannelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueServiceNowApprovalChannelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueServiceNowApprovalChannelResourceModel
	var details client.ApprovalChannelDetails
	details.Approver = &client.Approver{}
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	details.Approver.UserEmail = data.Approver.ValueString()
	details.Type = servicenow_approval_channel_type
	details.BaseUrl = data.BaseUrl.ValueStringPointer()
	details.UserName = data.UserName.ValueStringPointer()
	details.Password = data.Password.ValueStringPointer()
	details.Headers = data.Headers.ValueStringPointer()
	err := r.client.CreateApprovalChannel(data.Name.ValueString(), data.Description.ValueString(), details)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Approval Channel, got error: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueServiceNowApprovalChannelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueServiceNowApprovalChannelResourceModel

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
	data.BaseUrl = types.StringPointerValue(approval_channel.Details.BaseUrl)

	data.Headers = types.StringPointerValue(approval_channel.Details.Headers)
	data.Approver = types.StringPointerValue(&approval_channel.Details.Approver.UserEmail)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TorqueServiceNowApprovalChannelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueServiceNowApprovalChannelResourceModel
	var details client.ApprovalChannelDetails
	details.Approver = &client.Approver{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	details.Approver.UserEmail = data.Approver.ValueString()
	details.Type = servicenow_approval_channel_type
	details.BaseUrl = data.BaseUrl.ValueStringPointer()
	details.UserName = data.UserName.ValueStringPointer()
	details.Password = data.Password.ValueStringPointer()
	details.Headers = data.Headers.ValueStringPointer()
	err := r.client.UpdateApprovalChannel(data.Name.ValueString(), data.Description.ValueString(), details)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update Approval Channel, got error: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueServiceNowApprovalChannelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueServiceNowApprovalChannelResourceModel

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

func (r *TorqueServiceNowApprovalChannelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
