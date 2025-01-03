package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TorqueIntrospectionResource{}
var _ resource.ResourceWithImportState = &TorqueIntrospectionResource{}

func NewTorqueIntrospectionResource() resource.Resource {
	return &TorqueIntrospectionResource{}
}

// TorqueIntrospectionResource defines the resource implementation.
type TorqueIntrospectionResource struct {
}

// TorqueIntrospectionResourceModel describes the resource data model.
type TorqueIntrospectionResourceModel struct {
	DisplayName       types.String `tfsdk:"display_name"`
	Image             types.String `tfsdk:"image"`
	IntrospectionData types.Map    `tfsdk:"introspection_data"`
	Links             types.List   `tfsdk:"links"`
}

func (r *TorqueIntrospectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_introspection_resource"
}

func (r *TorqueIntrospectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Resource that will be presented in Torque resource catalog",

		Attributes: map[string]schema.Attribute{
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The resource name to show in Torque resource catalog",
				Required:            true,
			},
			"image": schema.StringAttribute{
				MarkdownDescription: "A link to an image for the custom resource. Can be hosted only on the following domains: `*.githubusercontent.com`, `*.quali.com`, `*.cloudfront.net`  ",
				Optional:            true,
				Computed:            false,
			},
			"introspection_data": schema.MapAttribute{
				MarkdownDescription: "Resource attribute to show in resource card. Note that only the first 4 attributes will be presented",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            false,
			},
			"links": schema.ListNestedAttribute{
				Description: "List of links that will be available as buttons in the resource introspection card.",
				Required:    false,
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"icon": schema.StringAttribute{
							Description: "Button's icon. Can be only one of the following: connect, restart, play, pause, stop, download, upload",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOf([]string{"connect", "restart", "play", "pause", "stop", "download", "upload"}...),
							},
						},
						"href": schema.StringAttribute{
							Description: "Button's link",
							Required:    true,
						},
						"label": schema.StringAttribute{
							Description: "Description that will be shown on hover",
							Required:    true,
						},
						"color": schema.StringAttribute{
							Description: "Hex value for the link's color",
							Optional:    true,
							Required:    false,
						},
					},
				},
			},
		},
	}
}

func (r *TorqueIntrospectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

}

func (r *TorqueIntrospectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueIntrospectionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	//     return
	// }

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	//data.Id = types.StringValue("example-id")

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueIntrospectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueIntrospectionResourceModel

	// Read Terraform prior state data into the model
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

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueIntrospectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueIntrospectionResourceModel

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

func (r *TorqueIntrospectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueIntrospectionResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }.
}

func (r *TorqueIntrospectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("display_name"), req, resp)
}
