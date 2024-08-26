package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/qualitorque/terraform-provider-torque/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TorqueAwsCostTargetResource{}
var _ resource.ResourceWithImportState = &TorqueAwsCostTargetResource{}

func NewTorqueAwsCostTargetResource() resource.Resource {
	return &TorqueAwsCostTargetResource{}
}

// TorqueAwsCostTargetResource defines the resource implementation.
type TorqueAwsCostTargetResource struct {
	client *client.Client
}

// TorqueAwsCostTargetResourceModel describes the resource data model.
type TorqueAwsCostTargetResourceModel struct {
	Name       types.String `tfsdk:"name"`
	RoleArn    types.String `tfsdk:"role_arn"`
	ExternalId types.String `tfsdk:"external_id"`
}

func (r *TorqueAwsCostTargetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_aws_cost_target"
}

func (r *TorqueAwsCostTargetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Creation of a new AWS Cost Collection target in Torque account.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the new cost collection target to be added to torque",
				Required:            true,
			},
			"role_arn": schema.StringAttribute{
				MarkdownDescription: "Please supply ARN Role with granted permission to query the AWS Cost Explorer API",
				Required:            true,
				Computed:            false,
			},
			"external_id": schema.StringAttribute{
				MarkdownDescription: "The AWS external id. For more infrormation: https://docs.qtorque.io/governance/cost-tracking/configuring-cost-aws",
				Required:            true,
				Computed:            false,
			},
		},
	}
}

func (r *TorqueAwsCostTargetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueAwsCostTargetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueAwsCostTargetResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.AddAWSCostTarget(data.Name.ValueString(), "aws", data.RoleArn.ValueString(), data.ExternalId.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create AWS cost collection target, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Resource Created Successful!")

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueAwsCostTargetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueAwsCostTargetResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueAwsCostTargetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state TorqueAwsCostTargetResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	current_target_name := state.Name
	err := r.client.UpdateAWSCostTarget(current_target_name.ValueString(), data.Name.ValueString(), "aws", data.RoleArn.ValueString(), data.ExternalId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating AWS Cost Target",
			"Could not update AWS Cost Target, unexpected error: "+err.Error(),
		)
		return
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueAwsCostTargetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueAwsCostTargetResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the space.
	err := r.client.DeleteCostTarget(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete AWS cost collection target, got error: %s", err))
		return
	}

}

func (r *TorqueAwsCostTargetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
