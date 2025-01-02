package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/qualitorque/terraform-provider-torque/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TorqueSpaceCustomIconResource{}
var _ resource.ResourceWithImportState = &TorqueSpaceCustomIconResource{}

func NewTorqueSpaceCustomIconResource() resource.Resource {
	return &TorqueSpaceCustomIconResource{}
}

// TorqueSpaceCustomIconResource defines the resource implementation.
type TorqueSpaceCustomIconResource struct {
	client *client.Client
}

// TorqueSpaceCustomIconResourceModel describes the resource data model.
type TorqueSpaceCustomIconResourceModel struct {
	SpaceName types.String `tfsdk:"space_name"`
	FilePath  types.String `tfsdk:"file_path"`
	FileName  types.String `tfsdk:"file_name"`
	Key       types.String `tfsdk:"key"`
}

func (r *TorqueSpaceCustomIconResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_space_custom_icon"
}

func (r *TorqueSpaceCustomIconResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Uploads an SVG file to be used as a custom icon for catalog items in Torque space.",

		Attributes: map[string]schema.Attribute{
			"space_name": schema.StringAttribute{
				MarkdownDescription: "Name of the space where this custom icon should be uploaded to Torque",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"file_path": schema.StringAttribute{
				MarkdownDescription: "Path to custom icon file. Supported format is SVG.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"file_name": schema.StringAttribute{
				MarkdownDescription: "Icon SVG file name.",
				Computed:            true,
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "Identifier for the icon, to be used in the catalog item resource when associating this icon with a catalog item.",
				Computed:            true,
			},
		},
	}
}

func (r *TorqueSpaceCustomIconResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueSpaceCustomIconResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueSpaceCustomIconResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UploadCustomIcon(data.SpaceName.ValueString(), data.FilePath.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to upload custom icon, got error: %s", err))
		return
	}
	var icon *client.TorqueSpaceCustomIcon
	icon, err = r.client.GetCustomIcon(data.SpaceName.ValueString(), data.FilePath.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read custom icon, got error: %s", err))
		return
	}
	data.FileName = types.StringPointerValue(&icon.FileName)
	data.Key = types.StringPointerValue(&icon.Key)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueSpaceCustomIconResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueSpaceCustomIconResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueSpaceCustomIconResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueSpaceCustomIconResourceModel
	diags := resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TorqueSpaceCustomIconResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueSpaceCustomIconResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the custom icon.
	err := r.client.DeleteCustomIcon(data.SpaceName.ValueString(), data.Key.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete custom icon, got error: %s", err))
		return
	}

}

func (r *TorqueSpaceCustomIconResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
