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
var _ resource.Resource = &TorqueSpaceLabelAssociationResource{}
var _ resource.ResourceWithImportState = &TorqueSpaceLabelAssociationResource{}

func NewTorqueSpaceLabelAssociationResource() resource.Resource {
	return &TorqueSpaceLabelAssociationResource{}
}

// TorqueSpaceLabelAssociationResource defines the resource implementation.
type TorqueSpaceLabelAssociationResource struct {
	client *client.Client
}

type TorqueSpaceLabelAssociationResourceModel struct {
	SpaceName      types.String `tfsdk:"space_name"`
	BlueprintName  types.String `tfsdk:"blueprint_name"`
	RepositoryName types.String `tfsdk:"repository_name"`
	Labels         types.List   `tfsdk:"labels"`
}

func (r *TorqueSpaceLabelAssociationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_space_label_association"
}

func (r *TorqueSpaceLabelAssociationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Associate Torque space label with a published blueprint (catalog item)",

		Attributes: map[string]schema.Attribute{
			"space_name": schema.StringAttribute{
				MarkdownDescription: "The name of the space where the catalog item and labels exist",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"blueprint_name": schema.StringAttribute{
				Description: "The name of the blueprint to associate the labels with",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"repository_name": schema.StringAttribute{
				Description: "The repository the blueprint is from",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"labels": schema.ListAttribute{
				Description: "List of label names to associate with the catalog item",
				Required:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *TorqueSpaceLabelAssociationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueSpaceLabelAssociationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueSpaceLabelAssociationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	labels := []string{}
	if !data.Labels.IsNull() {
		for _, label := range data.Labels.Elements() {
			labels = append(labels, strings.Trim(label.String(), "\""))
		}
	}
	err := r.client.EditCatalogItemLabels(data.SpaceName.ValueString(), data.BlueprintName.ValueString(), data.RepositoryName.ValueString(), labels)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to associate catalog item with label, got error: %s", err))
		return
	}

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueSpaceLabelAssociationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueSpaceLabelAssociationResourceModel

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

func (r *TorqueSpaceLabelAssociationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueSpaceLabelAssociationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	labels := []string{}
	if !data.Labels.IsNull() {
		for _, label := range data.Labels.Elements() {
			labels = append(labels, strings.Trim(label.String(), "\""))
		}
	}
	err := r.client.EditCatalogItemLabels(data.SpaceName.ValueString(), data.BlueprintName.ValueString(), data.RepositoryName.ValueString(), labels)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update catalog item label association, got error: %s", err))
		return
	}

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueSpaceLabelAssociationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueSpaceLabelAssociationResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	labels := []string{}
	err := r.client.EditCatalogItemLabels(data.SpaceName.ValueString(), data.BlueprintName.ValueString(), data.RepositoryName.ValueString(), labels)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update catalog item label association, got error: %s", err))
		return
	}

}

func (r *TorqueSpaceLabelAssociationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
