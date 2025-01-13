package resources

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/qualitorque/terraform-provider-torque/client"
	"github.com/qualitorque/terraform-provider-torque/internal/provider/common"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TorqueCatalogItemResource{}
var _ resource.ResourceWithImportState = &TorqueCatalogItemResource{}

func NewTorqueCatalogItemResource() resource.Resource {
	return &TorqueCatalogItemResource{}
}

// TorqueCatalogItemResource defines the resource implementation.
type TorqueCatalogItemResource struct {
	client *client.Client
}

// TorqueCatalogItemResourceModel describes the resource data model.
type TorqueCatalogItemResourceModel struct {
	SpaceName             types.String               `tfsdk:"space_name"`
	BlueprintName         types.String               `tfsdk:"blueprint_name"`
	DisplayName           types.String               `tfsdk:"display_name"`
	RepositoryName        types.String               `tfsdk:"repository_name"`
	MaxDuration           types.String               `tfsdk:"max_duration"`
	DefaultDuration       types.String               `tfsdk:"default_duration"`
	DefaultExtend         types.String               `tfsdk:"default_extend"`
	MaxActiveEnvironments types.Int32                `tfsdk:"max_active_environments"`
	AlwaysOn              types.Bool                 `tfsdk:"always_on"`
	AllowScheduling       types.Bool                 `tfsdk:"allow_scheduling"`
	CustomIcon            types.String               `tfsdk:"custom_icon"`
	Labels                types.List                 `tfsdk:"labels"`
	Tags                  []common.KeyValuePairModel `tfsdk:"tags"`
}

func (r *TorqueCatalogItemResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_catalog_item"
}

func (r *TorqueCatalogItemResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Creation of a new catalog item by publishing an existing blueprint to the self-service catalog.",

		Attributes: map[string]schema.Attribute{
			"space_name": schema.StringAttribute{
				MarkdownDescription: "Name of the space to configure",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"blueprint_name": schema.StringAttribute{
				MarkdownDescription: "The name of the blueprint to publish in the catalog",
				Required:            true,
				Computed:            false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the blueprint as it will be displayed in the self-service catalog.",
				Required:            false,
				Computed:            false,
				Optional:            true,
			},
			"repository_name": schema.StringAttribute{
				MarkdownDescription: "The name of the repository where the blueprint resides. \"Stored in Torque\" will be stored in \"qtorque\" repository",
				Required:            true,
				Computed:            false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"max_duration": schema.StringAttribute{
				MarkdownDescription: "The maximum duration of an environment instantiated from this blueprint.",
				Required:            false,
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^P(?:(\d+)D)?T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?$`),
						"must be a valid ISO 8601 timestamp (e.g., 2023-08-19T14:23:30Z or 2023-08-19T14:23:30+02:00)",
					),
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("always_on"),
					}...,
					),
				},
			},
			"default_duration": schema.StringAttribute{
				MarkdownDescription: "The default duration of an environment instantiated from this blueprint.",
				Required:            false,
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^P(?:(\d+)D)?T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?$`),
						"must be a valid ISO 8601 timestamp (e.g., 2023-08-19T14:23:30Z or 2023-08-19T14:23:30+02:00)",
					),
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("always_on"),
					}...),
				},
			},
			"default_extend": schema.StringAttribute{
				MarkdownDescription: "The default duration it will be possible to extend an environment instantiated from this blueprint.",
				Required:            false,
				Optional:            true,
				Computed:            true,

				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^P(?:(\d+)D)?T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?$`),
						"must be a valid ISO 8601 timestamp (e.g., 2023-08-19T14:23:30Z or 2023-08-19T14:23:30+02:00)",
					),
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("always_on"),
					}...),
				},
			},
			"max_active_environments": schema.Int32Attribute{
				MarkdownDescription: "Sets the maximum number of concurrent active environments insantiated from this blueprint.",
				Required:            false,
				Optional:            true,
			},
			"always_on": schema.BoolAttribute{
				MarkdownDescription: "Specify if environments launched from this blueprint should be always on or not.",
				Required:            false,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(
						path.Expressions{
							path.MatchRoot("default_duration"),
							path.MatchRoot("default_extend"),
							path.MatchRoot("max_duration"),
						}...,
					),
				},
			},
			"allow_scheduling": schema.BoolAttribute{
				MarkdownDescription: "Specify if environments from this blueprint can be scheduled to launch at a future time.",
				Required:            false,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"custom_icon": schema.StringAttribute{
				MarkdownDescription: "Custom icon key to associate with this catalog item. The key can be referenced from a torque_space_custom_icon key attribute.",
				Required:            false,
				Optional:            true,
				Computed:            false,
			},
			"labels": schema.ListAttribute{
				MarkdownDescription: "List of labels to associate with this catalog item.",
				Required:            false,
				Optional:            true,
				Computed:            false,
				ElementType:         types.StringType,
			},
			"tags": schema.ListNestedAttribute{
				Description: "Environment Tags",
				Required:    false,
				Optional:    true,
				Computed:    false,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Tag's name",
							Required:    true,
						},
						"value": schema.StringAttribute{
							Description: "The value of the tag",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func (r *TorqueCatalogItemResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueCatalogItemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueCatalogItemResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	if data.AlwaysOn.ValueBool() {
		data.MaxDuration = types.StringNull()
		data.DefaultExtend = types.StringNull()
		data.DefaultDuration = types.StringNull()
	} else {
		if data.MaxDuration.IsNull() || data.MaxDuration.IsUnknown() {
			data.MaxDuration = types.StringValue("PT2H")
		}
		if data.DefaultExtend.IsNull() || data.DefaultExtend.IsUnknown() {
			data.DefaultExtend = types.StringValue("PT2H")
		}
		if data.DefaultDuration.IsNull() || data.DefaultDuration.IsUnknown() {
			data.DefaultDuration = types.StringValue("PT2H")
		}
	}
	var maxActiveEnvironments *int32
	if !data.MaxActiveEnvironments.IsNull() {
		value := data.MaxActiveEnvironments.ValueInt32()
		maxActiveEnvironments = &value
	}
	if !data.CustomIcon.IsNull() {
		err := r.client.SetCatalogItemCustomIcon(data.SpaceName.ValueString(), data.BlueprintName.ValueString(), data.RepositoryName.ValueString(), data.CustomIcon.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Catalog Item, failed to set catalog item custom icon, got error: %s", err))
			return
		}
	}
	err := r.client.SetBlueprintPolicies(data.SpaceName.ValueString(), data.RepositoryName.ValueString(), data.BlueprintName.ValueString(), data.MaxDuration.ValueString(), data.DefaultDuration.ValueString(), data.DefaultExtend.ValueString(), maxActiveEnvironments, data.AlwaysOn.ValueBool(), data.AllowScheduling.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Catalog Item, failed to set blueprint policies, got error: %s", err))
		return
	}
	if !data.DisplayName.IsNull() {
		err = r.client.UpdateBlueprintDisplayName(data.SpaceName.ValueString(), data.RepositoryName.ValueString(), data.BlueprintName.ValueString(), data.DisplayName.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Catalog Item, failed to set blueprint display name, got error: %s", err))
			return
		}
	}
	err = r.client.PublishBlueprintInSpace(data.SpaceName.ValueString(), data.RepositoryName.ValueString(), data.BlueprintName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Catalog Item, failed to publish blueprint in space, got error: %s", err))
		return
	}
	if !data.Labels.IsNull() {
		var labels []string
		data.Labels.ElementsAs(ctx, &labels, false)
		err = r.client.EditCatalogItemLabels(data.SpaceName.ValueString(), data.BlueprintName.ValueString(), data.RepositoryName.ValueString(), labels)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Catalog Item, failed to associate labels, got error: %s", err))
			return
		}
	}
	for _, tag := range data.Tags {
		err = r.client.CreateBlueprintTagValue(data.SpaceName.ValueString(), tag.Name.ValueString(), tag.Value.ValueString(), data.RepositoryName.ValueString(), data.BlueprintName.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Catalog Item, failed to set blueprint tags, got error: %s", err))
			return
		}
	}
	tflog.Trace(ctx, "Resource Created Successful!")

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *TorqueCatalogItemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueCatalogItemResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	start := time.Now()
	for time.Since(start) < 20*time.Second {
		blueprint, err := r.client.GetBlueprint(data.SpaceName.ValueString(), data.BlueprintName.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get catalog items in space, got error: %s", err.Error()))
			return
		}

		if blueprint != nil && blueprint.Published {
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
			return
		}
		time.Sleep(500 * time.Millisecond) // Retry every 500ms
	}

	// Save updated data into Terraform state.
	resp.Diagnostics.AddError("Blueprint not published to catalog", "Blueprint was found in space but it is not published to catalog")
}

func (r *TorqueCatalogItemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueCatalogItemResourceModel
	var state TorqueCatalogItemResourceModel
	const default_icon = "nodes"

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.AlwaysOn.ValueBool() {
		data.MaxDuration = types.StringNull()
		data.DefaultExtend = types.StringNull()
		data.DefaultDuration = types.StringNull()
	} else {
		if data.MaxDuration.IsNull() || data.MaxDuration.IsUnknown() {
			data.MaxDuration = types.StringValue("PT2H")
		}
		if data.DefaultExtend.IsNull() || data.DefaultExtend.IsUnknown() {
			data.DefaultExtend = types.StringValue("PT2H")
		}
		if data.DefaultDuration.IsNull() || data.DefaultDuration.IsUnknown() {
			data.DefaultDuration = types.StringValue("PT2H")
		}
	}

	var maxActiveEnvironments *int32
	if !data.MaxActiveEnvironments.IsNull() {
		value := data.MaxActiveEnvironments.ValueInt32()
		maxActiveEnvironments = &value
	}

	err := r.client.SetBlueprintPolicies(data.SpaceName.ValueString(), data.RepositoryName.ValueString(), data.BlueprintName.ValueString(), data.MaxDuration.ValueString(), data.DefaultDuration.ValueString(), data.DefaultExtend.ValueString(), maxActiveEnvironments, data.AlwaysOn.ValueBool(), data.AllowScheduling.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to set blueprint policies, got error: %s", err))
		return
	}
	if !data.DisplayName.IsNull() {
		err = r.client.UpdateBlueprintDisplayName(data.SpaceName.ValueString(), data.RepositoryName.ValueString(), data.BlueprintName.ValueString(), data.DisplayName.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update Catalog Item, failed to set blueprint display name, got error: %s", err))
			return
		}
	}
	if data.CustomIcon.IsNull() {
		err := r.client.SetCatalogItemIcon(data.SpaceName.ValueString(), data.BlueprintName.ValueString(), data.RepositoryName.ValueString(), default_icon)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to remove Catalog Item custom icon, failed to set catalog item custom icon, got error: %s", err))
			return
		}
	} else {
		if data.CustomIcon.ValueString() != state.CustomIcon.ValueString() {
			err := r.client.SetCatalogItemCustomIcon(data.SpaceName.ValueString(), data.BlueprintName.ValueString(), data.RepositoryName.ValueString(), data.CustomIcon.ValueString())
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update workflow custom icon, failed to set catalog item custom icon, got error: %s", err))
				return
			}
		}
	}

	for _, tag := range data.Tags {
		err = r.client.SetBlueprintTagValue(data.SpaceName.ValueString(), tag.Name.ValueString(), tag.Value.ValueString(), data.RepositoryName.ValueString(), data.BlueprintName.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Catalog Item, failed to set blueprint tags, got error: %s", err))
			return
		}
	}

	var labels []string
	data.Labels.ElementsAs(ctx, &labels, false)
	err = r.client.EditCatalogItemLabels(data.SpaceName.ValueString(), data.BlueprintName.ValueString(), data.RepositoryName.ValueString(), labels)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update Catalog Item, failed to update labels, got error: %s", err))
		return
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueCatalogItemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueCatalogItemResourceModel
	const default_icon = "nodes"
	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UnpublishBlueprintInSpace(data.SpaceName.ValueString(), data.RepositoryName.ValueString(), data.BlueprintName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to unpublish blueprint from space, got error: %s", err))
		return
	}
	err = r.client.SetCatalogItemIcon(data.SpaceName.ValueString(), data.BlueprintName.ValueString(), data.RepositoryName.ValueString(), default_icon)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to remove Catalog Item custom icon, failed to set catalog item custom icon, got error: %s", err))
		return
	}

	data.MaxDuration = types.StringValue("PT2H")
	data.DefaultExtend = types.StringValue("PT2H")
	data.DefaultDuration = types.StringValue("PT2H")

	err = r.client.SetBlueprintPolicies(data.SpaceName.ValueString(), data.RepositoryName.ValueString(), data.BlueprintName.ValueString(), data.MaxDuration.ValueString(), data.DefaultDuration.ValueString(), data.DefaultExtend.ValueString(), nil, false, false)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to set blueprint policies, got error: %s", err))
		return
	}
	err = r.client.UpdateBlueprintDisplayName(data.SpaceName.ValueString(), data.RepositoryName.ValueString(), data.BlueprintName.ValueString(), data.BlueprintName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Catalog Item, failed to set blueprint display name, got error: %s", err))
		return
	}
	for _, tag := range data.Tags {
		err = r.client.DeleteBlueprintTagValue(data.SpaceName.ValueString(), tag.Name.ValueString(), data.RepositoryName.ValueString(), data.BlueprintName.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Catalog Item, failed to delete blueprint tags, got error: %s", err))
			return
		}
	}
	err = r.client.EditCatalogItemLabels(data.SpaceName.ValueString(), data.BlueprintName.ValueString(), data.RepositoryName.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to remove Catalog Item custom icon, failed to set catalog item custom icon, got error: %s", err))
		return
	}
}

func (r *TorqueCatalogItemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
