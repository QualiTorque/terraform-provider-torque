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
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/qualitorque/terraform-provider-torque/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TorqueEnvironmentLabelAssociationResource{}
var _ resource.ResourceWithImportState = &TorqueEnvironmentLabelAssociationResource{}

func NewTorqueEnvironmentLabelAssociationResource() resource.Resource {
	return &TorqueEnvironmentLabelAssociationResource{}
}

// TorqueEnvironmentLabelAssociationResource defines the resource implementation.
type TorqueEnvironmentLabelAssociationResource struct {
	client *client.Client
}

// TorqueEnvironmentLabelAssociationResourceModel describes the resource data model.
type TorqueEnvironmentLabelAssociationResourceModel struct {
	SpaceName     types.String        `tfsdk:"space_name"`
	EnvironmentId types.String        `tfsdk:"environment_id"`
	Labels        []keyValuePairModel `tfsdk:"labels"`
}

type keyValuePairModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

func (r *TorqueEnvironmentLabelAssociationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_environment_label_association"
}

func (r *TorqueEnvironmentLabelAssociationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Creation of a new label to be used in Torque environment and can be associated to catalog items.",

		Attributes: map[string]schema.Attribute{
			"space_name": schema.StringAttribute{
				MarkdownDescription: "Space where the environment resides in",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Environment id to associate the labels with",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"labels": schema.ListNestedAttribute{
				Description: "List of labels associated with the environment.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Description: "Input's name",
							Required:    true,
						},
						"value": schema.StringAttribute{
							Description: "Input's default value",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func (r *TorqueEnvironmentLabelAssociationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueEnvironmentLabelAssociationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueEnvironmentLabelAssociationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	added_labels := []client.KeyValuePair{}
	if len(data.Labels) > 0 {
		for _, label := range data.Labels {

			added_labels = append(added_labels, client.KeyValuePair{
				Key:   label.Key.ValueString(),
				Value: label.Value.ValueString(),
			})
		}
	}
	
	err := r.client.UpdateEnvironmentLabels(data.EnvironmentId.ValueString(), data.SpaceName.ValueString(), added_labels, nil)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create environment label, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Resource Created Successful!")

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueEnvironmentLabelAssociationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueEnvironmentLabelAssociationResourceModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var env_labels []keyValuePairModel
	labels, err := r.client.GetEnvironmentLabels(data.SpaceName.ValueString(), data.EnvironmentId.ValueString())
	if len(labels) > 0 {
		for _, label := range labels {
			env_labels = append(env_labels, keyValuePairModel{
				Key:   types.StringValue(label.Key),
				Value: types.StringValue(label.Value),
			})
		}
	}
	// label, err := r.client.GetEnvironmentLabel(data.Key.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading environment labels details",
			"Could not read Torque label name "+data.EnvironmentId.ValueString()+": "+err.Error(),
		)
		return
	}

	// Treat HTTP 404 Not Found status as a signal to recreate resource
	// and return early
	// if label.Value == "" {
	// 	tflog.Error(ctx, "label not found in Torque")
	// 	resp.State.RemoveResource(ctx)
	// 	return
	// }
	// data.Key = types.StringValue(label.Key)
	// data.Value = types.StringValue(label.Value)
	// Set refreshed state
	data.Labels = env_labels
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TorqueEnvironmentLabelAssociationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state TorqueEnvironmentLabelAssociationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateLabels := state.Labels
	planLabels := data.Labels

	added_labels := []client.KeyValuePair{}
	if len(planLabels) > 0 {
		for _, label := range planLabels {
			if !labelExists(label, stateLabels) {
				added_labels = append(added_labels, client.KeyValuePair{
					Key:   label.Key.ValueString(),
					Value: label.Value.ValueString(),
				})
			}
		}
	}
	removed_labels := []client.KeyValuePair{}
	if len(stateLabels) > 0 {
		for _, label := range stateLabels {
			if !labelExists(label, planLabels) {
				removed_labels = append(removed_labels, client.KeyValuePair{
					Key:   label.Key.ValueString(),
					Value: label.Value.ValueString(),
				})
			}
		}
	}
	// if len(data.Labels) > 0 {
	// 	for _, label := range data.Labels {

	// 		removed_labels = append(removed_labels, client.KeyValuePair{
	// 			Key:   label.Key.ValueString(),
	// 			Value: label.Value.ValueString(),
	// 		})
	// 	}
	// }

	err := r.client.UpdateEnvironmentLabels(data.EnvironmentId.ValueString(), data.SpaceName.ValueString(), added_labels, removed_labels)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update environment labels, got error: %s", err))
		return
	}
	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TorqueEnvironmentLabelAssociationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueEnvironmentLabelAssociationResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	removed_labels := []client.KeyValuePair{}
	if len(data.Labels) > 0 {
		for _, label := range data.Labels {

			removed_labels = append(removed_labels, client.KeyValuePair{
				Key:   label.Key.ValueString(),
				Value: label.Value.ValueString(),
			})
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.UpdateEnvironmentLabels(data.EnvironmentId.ValueString(), data.SpaceName.ValueString(), nil, removed_labels)

	// Delete the environment label.
	// err := r.client.DeleteEnvironmentLabel(data.Key.ValueString(), data.Value.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete environment label, got error: %s", err))
		return
	}

}

func (r *TorqueEnvironmentLabelAssociationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func labelExists(label keyValuePairModel, labels []keyValuePairModel) bool {
	for _, l := range labels {
		if l.Key == label.Key && l.Value == label.Value {
			return true
		}
	}
	return false
}
