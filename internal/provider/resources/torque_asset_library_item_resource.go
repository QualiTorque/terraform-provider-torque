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
	"github.com/qualitorque/terraform-provider-torque/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TorqueAssetLibraryItemResource{}
var _ resource.ResourceWithImportState = &TorqueAssetLibraryItemResource{}

func NewTorqueAssetLibraryItemResource() resource.Resource {
	return &TorqueAssetLibraryItemResource{}
}

// TorqueAssetLibraryItemResource defines the resource implementation.
type TorqueAssetLibraryItemResource struct {
	client *client.Client
}

// torqueAssetLibraryItemResource describes the resource data model.
type TorqueAssetLibraryItemResourceModel struct {
	SpaceName      types.String `tfsdk:"space_name"`
	BlueprintName  types.String `tfsdk:"blueprint_name"`
	RepositoryName types.String `tfsdk:"repository_name"`
}

func (r *TorqueAssetLibraryItemResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_asset_library_item"
}

func (r *TorqueAssetLibraryItemResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Adds a blueprint to an asset-library so it can serve as a building block.",

		Attributes: map[string]schema.Attribute{
			"space_name": schema.StringAttribute{
				MarkdownDescription: "Name of the space where the blueprint is",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"blueprint_name": schema.StringAttribute{
				MarkdownDescription: "The name of the blueprint to add to the asset-library",
				Required:            true,
				Computed:            false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"repository_name": schema.StringAttribute{
				MarkdownDescription: "The name of the repository where the blueprint resides. \"Stored in Torque\" will be stored in \"qtorque\" repository",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *TorqueAssetLibraryItemResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueAssetLibraryItemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueAssetLibraryItemResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.AddBlueprintToAssetLibrary(data.SpaceName.ValueString(), data.RepositoryName.ValueStringPointer(), data.BlueprintName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to add blueprint to asset-library, got error: %s", err))
		return
	}

	if data.RepositoryName.IsUnknown() {
		blueprint, err := r.client.GetBlueprint(data.SpaceName.ValueString(), data.BlueprintName.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to retrieve blueprint details, got error: %s", err))
			return
		}
		data.RepositoryName = types.StringValue(blueprint.RepoName)

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *TorqueAssetLibraryItemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TorqueAssetLibraryItemResourceModel

	// Read Terraform prior state data into the model.
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	_, err := r.client.GetBlueprintFromAssetLibrary(state.SpaceName.ValueString(), state.BlueprintName.ValueString())
	if err != nil {
		// Check if the error is a NotFoundError and remove the resource from state
		if strings.Contains(err.Error(), "not found") {
			resp.Diagnostics.AddWarning("Asset-Library item will be recreated","Blueprint was removed from asset-library outside of Terraform.")
			resp.State.RemoveResource(ctx)
			return
		}
		// Otherwise, return the error
		resp.Diagnostics.AddError(
			"Error reading blueprint",
			fmt.Sprintf("Could not read blueprint %s in space %s: %s", state.BlueprintName.ValueString(), state.SpaceName.ValueString(), err.Error()),
		)
		return
	}

	// Save updated data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *TorqueAssetLibraryItemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueAssetLibraryItemResourceModel

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

func (r *TorqueAssetLibraryItemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueAssetLibraryItemResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the space.
	err := r.client.RemoveBlueprintFromAssetLibrary(data.SpaceName.ValueString(), data.RepositoryName.ValueStringPointer(), data.BlueprintName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to remove blueprint from asset-library, got error: %s", err))
		return
	}

}

func (r *TorqueAssetLibraryItemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
