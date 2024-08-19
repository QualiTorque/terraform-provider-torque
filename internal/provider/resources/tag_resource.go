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
var _ resource.Resource = &TorqueTagResource{}
var _ resource.ResourceWithImportState = &TorqueTagResource{}

func NewTorqueTagResource() resource.Resource {
	return &TorqueTagResource{}
}

// TorqueTagResource defines the resource implementation.
type TorqueTagResource struct {
	client *client.Client
}

// POST  https://portal.qtorque.io/api/settings/tags
// {"name":"tomera","value":"unassigned_tag_value","scope":"space","description":"asd","possible_values":["a","b"]}

// TorqueTagResourceModel describes the resource data model.
type TorqueTagResourceModel struct {
	Name           types.String `tfsdk:"name"`
	Value          types.String `tfsdk:"value"`
	Scope          types.String `tfsdk:"scope"`
	Description    types.String `tfsdk:"description"`
	PossibleValues types.List   `tfsdk:"possible_values"`
}

func (r *TorqueTagResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_tag"
}

func (r *TorqueTagResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Creation of a new Torque Tag, it's scope, value and/or possible values",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the new tag to be added to torque",
				Required:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "Tag value to be set as the tag value default",
				Required:            true,
				Computed:            false,
			},
			"scope": schema.StringAttribute{
				MarkdownDescription: "Tag scope. Possible values: account, space, blueprint, environment",
				Required:            true,
				Computed:            false,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Tag description",
				Optional:            true,
				Computed:            false,
			},
			"possible_values": schema.ListAttribute{
				MarkdownDescription: "Tag possible values",
				Optional:            true,
				Computed:            false,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *TorqueTagResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueTagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueTagResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var possible []string
	if !data.PossibleValues.IsNull() {
		for _, pos_value := range data.PossibleValues.Elements() {
			possible = append(possible, pos_value.String())
		}
	}

	err := r.client.AddTag(data.Name.ValueString(), data.Value.ValueString(), data.Description.ValueString(), possible, data.Scope.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create tag, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Resource Created Successful!")

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *TorqueTagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueTagResourceModel

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

func (r *TorqueTagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueTagResourceModel

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

func (r *TorqueTagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueTagResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the space.
	err := r.client.RemoveTag(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete tag, got error: %s", err))
		return
	}

}

func (r *TorqueTagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
