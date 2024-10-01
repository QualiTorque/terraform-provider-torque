package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/qualitorque/terraform-provider-torque/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TorqueTagSpaceValueAssociationResource{}
var _ resource.ResourceWithImportState = &TorqueTagSpaceValueAssociationResource{}

func NewTorqueTagSpaceValueAssociationResource() resource.Resource {
	return &TorqueTagSpaceValueAssociationResource{}
}

// TorqueTagSpaceValueAssociationResource defines the resource implementation.
type TorqueTagSpaceValueAssociationResource struct {
	client *client.Client
}

type TorqueTagSpaceValueAssociationResourceModel struct {
	SpaceName types.String `tfsdk:"space_name"`
	TagName   types.String `tfsdk:"tag_name"`
	TagValue  types.String `tfsdk:"tag_value"`
}

func (r *TorqueTagSpaceValueAssociationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_space_tag_value_association"
}

func (r *TorqueTagSpaceValueAssociationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Associate Torque space with existing tag and value",

		Attributes: map[string]schema.Attribute{
			"space_name": schema.StringAttribute{
				MarkdownDescription: "Existing Torque Space name",
				Required:            true,
			},
			"tag_name": schema.StringAttribute{
				Description: "Tag name configured in the account",
				Required:    true,
			},
			"tag_value": schema.StringAttribute{
				Description: "The tag value to be set for the space",
				Required:    true,
			},
		},
	}
}

func (r *TorqueTagSpaceValueAssociationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueTagSpaceValueAssociationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueTagSpaceValueAssociationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CreateSpaceTagValue(data.SpaceName.ValueString(), data.TagName.ValueString(), data.TagValue.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "422") {
			newErr := r.client.SetSpaceTagValue(data.SpaceName.ValueString(), data.TagName.ValueString(), data.TagValue.ValueString())
			if newErr != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to set tag value in space, got error: %s", newErr))
				return
			}
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create tag value in space, got error: %s", err))
			return
		}
	}
	tflog.Trace(ctx, "Resource Created Successful!")

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueTagSpaceValueAssociationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueTagSpaceValueAssociationResourceModel

	// Read Terraform prior state data into the model.
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tag, err := r.client.GetSpaceTag(data.SpaceName.ValueString(), data.TagName.ValueString())
	if tag == (client.NameValuePair{}) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading tag details bla bla",
			"Could not read space tag value of "+data.TagName.ValueString()+": "+err.Error(),
		)
		return
	}
	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }
	data.TagName = types.StringValue(tag.Name)
	data.TagValue = types.StringValue(tag.Value)

	// Save updated data into Terraform state.
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TorqueTagSpaceValueAssociationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueTagSpaceValueAssociationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.SetSpaceTagValue(data.SpaceName.ValueString(), data.TagName.ValueString(), data.TagValue.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to set tag value in space, got error: %s", err))
		return
	}
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

func (r *TorqueTagSpaceValueAssociationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueTagSpaceValueAssociationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSpaceTagValue(data.SpaceName.ValueString(), data.TagName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete tag value in space, got error: %s", err))
		return
	}
	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *TorqueTagSpaceValueAssociationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
