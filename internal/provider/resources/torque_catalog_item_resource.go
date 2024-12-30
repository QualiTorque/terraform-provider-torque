package resources

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/qualitorque/terraform-provider-torque/client"
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
	SpaceName             types.String `tfsdk:"space_name"`
	BlueprintName         types.String `tfsdk:"blueprint_name"`
	RepositoryName        types.String `tfsdk:"repository_name"`
	MaxDuration           types.String `tfsdk:"max_duration"`
	DefaultDuration       types.String `tfsdk:"default_duration"`
	DefaultExtend         types.String `tfsdk:"default_extend"`
	MaxActiveEnvironments types.Int32  `tfsdk:"max_active_environments"`
	AlwaysOn              types.Bool   `tfsdk:"always_on"`
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
				// Computed:            true,
				// Default:             stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})$`),
						"must be a valid ISO 8601 timestamp (e.g., 2023-08-19T14:23:30Z or 2023-08-19T14:23:30+02:00)",
					),
				},
			},
			"default_duration": schema.StringAttribute{
				MarkdownDescription: "The default duration of an environment instantiated from this blueprint.",
				Required:            false,
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("PT2H"),
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})$`),
						"must be a valid ISO 8601 timestamp (e.g., 2023-08-19T14:23:30Z or 2023-08-19T14:23:30+02:00)",
					),
				},
			},
			"default_extend": schema.StringAttribute{
				MarkdownDescription: "The default duration it will be possible to extend an environment instantiated from this blueprint.",
				Required:            false,
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("PT2H"),
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})$`),
						"must be a valid ISO 8601 timestamp (e.g., 2023-08-19T14:23:30Z or 2023-08-19T14:23:30+02:00)",
					),
				},
			},
			"max_active_environments": schema.Int32Attribute{
				MarkdownDescription: "Sets the maximum number of concurrent active environments insantiated from this blueprint.",
				Required:            false,
				Optional:            true,
				// Computed:            true,
				// PlanModifiers: []planmodifier.Int32{
				// 	int32planmodifier.UseStateForUnknown(),
				// },
				// Default: defaults.Int32.DefaultInt32(-1),
			},
			"always_on": schema.BoolAttribute{
				MarkdownDescription: "Specify if environments launched from this blueprint should be always on or not.",
				Required:            false,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
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
	var maxActiveEnvironments *int32
	if !data.MaxActiveEnvironments.IsNull() {
		value := data.MaxActiveEnvironments.ValueInt32() // Get the actual value
		maxActiveEnvironments = &value                   // Pass it as a pointer
	}
	err := r.client.SetBlueprintPolicies(data.SpaceName.ValueString(), data.RepositoryName.ValueString(), data.BlueprintName.ValueString(), data.MaxDuration.ValueString(), data.DefaultDuration.ValueString(), data.DefaultDuration.ValueString(), maxActiveEnvironments, data.AlwaysOn.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to set blueprint policies, got error: %s", err))
		return
	}

	err = r.client.PublishBlueprintInSpace(data.SpaceName.ValueString(), data.RepositoryName.ValueString(), data.BlueprintName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to publish blueprint in space, got error: %s", err))
		return
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

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	var maxActiveEnvironments *int32
	if !data.MaxActiveEnvironments.IsNull() {
		value := data.MaxActiveEnvironments.ValueInt32() // Get the actual value
		maxActiveEnvironments = &value                   // Pass it as a pointer
	}
	err := r.client.SetBlueprintPolicies(data.SpaceName.ValueString(), data.RepositoryName.ValueString(), data.BlueprintName.ValueString(), data.MaxDuration.ValueString(), data.DefaultDuration.ValueString(), data.DefaultDuration.ValueString(), maxActiveEnvironments, data.AlwaysOn.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to set blueprint policies, got error: %s", err))
		return
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueCatalogItemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueCatalogItemResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the space.
	err := r.client.UnpublishBlueprintInSpace(data.SpaceName.ValueString(), data.RepositoryName.ValueString(), data.BlueprintName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to unpublish blueprint from space, got error: %s", err))
		return
	}

}

func (r *TorqueCatalogItemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
