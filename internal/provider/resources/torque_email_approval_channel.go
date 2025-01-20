package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the approval channel",
				Optional:            true,
				Computed:            false,
				Required:            false,
			},
			"approvers": schema.ListAttribute{
				Description: "List of spaces that can use this approval channel",
				Required:    false,
				Optional:    true,
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

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// input_source, err := r.client.GetInputSource(data.Name.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error Reading Input Source details",
	// 		"Could not read Input Source "+data.Name.ValueString()+": "+err.Error(),
	// 	)
	// 	return
	// }
	// data.Name = types.StringValue(input_source.Name)
	// data.Description = types.StringValue(input_source.Description)
	// data.BucketName = types.StringValue(input_source.Details.BucketName.Value)
	// data.BucketNameOverridable = types.BoolValue(input_source.Details.BucketName.Overridable)
	// data.CredentialName = types.StringValue(input_source.Details.CredentialName)
	// data.AllSpaces = types.BoolValue(input_source.AllowedSpaces.AllSpaces)
	// if len(input_source.AllowedSpaces.SpecificSpaces) > 0 {
	// 	data.SpecificSpaces, _ = types.ListValueFrom(ctx, types.StringType, input_source.AllowedSpaces.SpecificSpaces)
	// } else {
	// 	data.SpecificSpaces = types.ListNull(types.StringType)
	// }
	// data.FilterPattern = types.StringValue(input_source.Details.FilterPattern.Value)
	// data.FilterPatternOverridable = types.BoolValue(input_source.Details.FilterPattern.Overridable)
	// data.PathPrefix = types.StringValue(input_source.Details.PathPrefix.Value)
	// data.PathPrefixOverridable = types.BoolValue(input_source.Details.PathPrefix.Overridable)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TorqueEmailApprovalChannelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueEmailApprovalChannelResourceModel
	// var state TorqueEmailApprovalChannelResourceModel
	// var details client.InputSourceDetails
	// var allowed_spaces client.AllowedSpaces
	// const input_source_type = "s3-object"
	// var specificSpaces []string
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	// resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	// if resp.Diagnostics.HasError() {
	// 	return
	// }
	// allowed_spaces.AllSpaces = data.AllSpaces.ValueBool()
	// if !data.SpecificSpaces.IsNull() {
	// 	allowed_spaces.AllSpaces = false
	// 	for _, val := range data.SpecificSpaces.Elements() {
	// 		specificSpaces = append(specificSpaces, strings.Replace(val.String(), "\"", "", -1))
	// 	}
	// 	allowed_spaces.SpecificSpaces = specificSpaces
	// } else {
	// 	allowed_spaces.AllSpaces = data.AllSpaces.ValueBool() // true
	// }
	// details.BucketName.Overridable = data.BucketNameOverridable.ValueBool()
	// details.BucketName.Value = data.BucketName.ValueString()
	// details.FilterPattern.Overridable = data.FilterPatternOverridable.ValueBool()
	// details.FilterPattern.Value = data.FilterPattern.ValueString()
	// details.PathPrefix.Overridable = data.PathPrefixOverridable.ValueBool()
	// details.PathPrefix.Value = data.PathPrefix.ValueString()
	// details.Type = input_source_type
	// details.CredentialName = data.CredentialName.ValueString()
	// err := r.client.UpdateInputSource(state.Name.ValueString(), data.Name.ValueString(), data.Description.ValueString(), allowed_spaces, details)
	// if err != nil {
	// 	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Input Source, got error: %s", err))
	// 	return
	// }
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
