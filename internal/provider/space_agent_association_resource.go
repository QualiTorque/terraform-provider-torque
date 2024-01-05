package provider

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
var _ resource.Resource = &TorqueAgentSpaceAssociationResource{}
var _ resource.ResourceWithImportState = &TorqueAgentSpaceAssociationResource{}

func NewTorqueAgentSpaceAssociationResource() resource.Resource {
	return &TorqueAgentSpaceAssociationResource{}
}

// TorqueAgentSpaceAssociationResource defines the resource implementation.
type TorqueAgentSpaceAssociationResource struct {
	client *client.Client
}

type TorqueAgentSpaceAssociationResourceModel struct {
	SpaceName      types.String `tfsdk:"space_name"`
	AgentName      types.String `tfsdk:"agent_name"`
	Namespace      types.String `tfsdk:"namespace"`
	ServiceAccount types.String `tfsdk:"service_account"`
}

func (r *TorqueAgentSpaceAssociationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_agent_space_association"
}

func (r *TorqueAgentSpaceAssociationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Associate Torque space with existing registered agent",

		Attributes: map[string]schema.Attribute{
			"space_name": schema.StringAttribute{
				MarkdownDescription: "Existing Torque Space name",
				Required:            true,
			},
			"agent_name": schema.StringAttribute{
				Description: "Agent name to associate to the space",
				Required:    true,
			},
			"service_account": schema.StringAttribute{
				Description: "Default service account to be used with the agent in the space",
				Required:    true,
			},
			"namespace": schema.StringAttribute{
				Description: "Default namespace to be used with the agent in the space",
				Required:    true,
			},
		},
	}
}

func (r *TorqueAgentSpaceAssociationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueAgentSpaceAssociationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueAgentSpaceAssociationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.AddAgentToSpace(data.AgentName.ValueString(), data.Namespace.ValueString(),
		data.ServiceAccount.ValueString(), data.SpaceName.ValueString(), "K8S")
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to attach agent to space, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Resource Created Successful!")

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueAgentSpaceAssociationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueAgentSpaceAssociationResourceModel

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

func (r *TorqueAgentSpaceAssociationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueAgentSpaceAssociationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueAgentSpaceAssociationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueAgentSpaceAssociationResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Remove repos from space.
	err := r.client.RemoveAgentFromSpace(data.AgentName.ValueString(), data.SpaceName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to attach agent to space, got error: %s", err))
		return
	}

}

func (r *TorqueAgentSpaceAssociationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
