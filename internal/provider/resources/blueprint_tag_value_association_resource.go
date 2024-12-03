package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/qualitorque/terraform-provider-torque/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TorqueTagBlueprintValueAssociationResource{}
var _ resource.ResourceWithImportState = &TorqueTagBlueprintValueAssociationResource{}

func NewTorqueTagBlueprintValueAssociationResource() resource.Resource {
	return &TorqueTagBlueprintValueAssociationResource{}
}

// TorqueTagBlueprintValueAssociationResource defines the resource implementation.
type TorqueTagBlueprintValueAssociationResource struct {
	client *client.Client
}

type TorqueTagBlueprintValueAssociationResourceModel struct {
	SpaceName      types.String `tfsdk:"space_name"`
	RepositoryName types.String `tfsdk:"repository_name"`
	TagName        types.String `tfsdk:"tag_name"`
	TagValue       types.String `tfsdk:"tag_value"`
	BlueprintName  types.String `tfsdk:"blueprint_name"`
}

func (r *TorqueTagBlueprintValueAssociationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_blueprint_tag_value_association"
}

func (r *TorqueTagBlueprintValueAssociationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Associate Torque tag value in a blueprint",

		Attributes: map[string]schema.Attribute{
			"space_name": schema.StringAttribute{
				MarkdownDescription: "Existing Torque Space name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tag_name": schema.StringAttribute{
				Description: "The Tag name configured at the account level with a 'blueprint' scope",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tag_value": schema.StringAttribute{
				Description: "The tag value to be set for the blueprint",
				Required:    true,
			},
			"blueprint_name": schema.StringAttribute{
				Description: "The blueprint to be used",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"repository_name": schema.StringAttribute{
				Description: "The blueprint repository where the blueprint is stored. for \"stored in Torque\" use 'qtorque'",
				Required:    true,
			},
		},
	}
}

func (r *TorqueTagBlueprintValueAssociationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueTagBlueprintValueAssociationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueTagBlueprintValueAssociationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CreateBlueprintTagValue(data.SpaceName.ValueString(), data.TagName.ValueString(),
		data.TagValue.ValueString(), data.RepositoryName.ValueString(), data.BlueprintName.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "422") {
			new_err := r.client.SetBlueprintTagValue(data.SpaceName.ValueString(), data.TagName.ValueString(), data.TagValue.ValueString(), data.RepositoryName.ValueString(), data.BlueprintName.ValueString())
			if new_err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to set tag value in blueprint, got error: %s", err))
				return
			}
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create tag value in blueprint, got error: %s", err))
			return
		}
	}
	tflog.Trace(ctx, "Resource Created Successful!")

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueTagBlueprintValueAssociationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueTagBlueprintValueAssociationResourceModel

	// Read Terraform prior state data into the model.
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tag, err := r.client.GetBlueprintTag(data.SpaceName.ValueString(), data.TagName.ValueString(), data.RepositoryName.ValueString(), data.BlueprintName.ValueString())
	if tag == (client.NameValuePair{}) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading tag details",
			"Could not read blueprint tag "+data.TagName.ValueString()+": "+err.Error(),
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

func (r *TorqueTagBlueprintValueAssociationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueTagBlueprintValueAssociationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.SetBlueprintTagValue(data.SpaceName.ValueString(), data.TagName.ValueString(), data.TagValue.ValueString(), data.RepositoryName.ValueString(), data.BlueprintName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to set tag value in space, got error: %s", err))
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

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TorqueTagBlueprintValueAssociationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueTagBlueprintValueAssociationResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.DeleteBlueprintTagValue(data.SpaceName.ValueString(), data.TagName.ValueString(), data.RepositoryName.ValueString(), data.BlueprintName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete tag value in space, got error: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TorqueTagBlueprintValueAssociationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
