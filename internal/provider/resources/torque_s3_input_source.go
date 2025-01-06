package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/qualitorque/terraform-provider-torque/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TorqueS3ObjectInputSourceResource{}
var _ resource.ResourceWithImportState = &TorqueS3ObjectInputSourceResource{}

func NewTorqueS3ObjectInputSourceResource() resource.Resource {
	return &TorqueS3ObjectInputSourceResource{}
}

// TorqueS3ObjectInputSourceResource defines the resource implementation.
type TorqueS3ObjectInputSourceResource struct {
	client *client.Client
}

// TorqueS3ObjectInputSourceResourceModel describes the resource data model.
type TorqueS3ObjectInputSourceResourceModel struct {
	Name                     types.String `tfsdk:"name"`
	Description              types.String `tfsdk:"description"`
	AllSpaces                types.Bool   `tfsdk:"all_spaces"`
	SpecificSpaces           types.List   `tfsdk:"specific_spaces"`
	BucketName               types.String `tfsdk:"bucket_name"`
	BucketNameOverridable    types.Bool   `tfsdk:"bucket_name_overridable"`
	CredentialName           types.String `tfsdk:"credential_name"`
	FilterPattern            types.String `tfsdk:"filter_pattern"`
	FilterPatternOverridable types.Bool   `tfsdk:"filter_pattern_overridable"`
	PathPrefix               types.String `tfsdk:"path_prefix"`
	PathPrefixOverridable    types.Bool   `tfsdk:"path_prefix_overridable"`
	// Type                     types.String `tfsdk:"type"`
}

// type TorqueInputSource struct {
// 	Name          string             `json:"file_name"`
// 	Description   string             `json:"key"`
// 	Details       InputSourceDetails `json:"details"`
// 	AllowedSpaces AllowedSpaces      `json:"allowed_spaces"`
// }

type overridableValue struct {
	Overridable types.Bool   `tfsdk:"overridable"`
	Value       types.String `tfsdk:"value"`
}

type allowedSpaces struct {
	AllSpaces      types.Bool `tfsdk:"all_spaces"`
	SpecificSpaces types.List `tfsdk:"specific_spaces"`
}

func (r *TorqueS3ObjectInputSourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_s3_object_input_source"
}

func (r *TorqueS3ObjectInputSourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Allows to enable and publish existing Torque workflow with env or env_resource scope so it will be allowed to be executed and displayed in the self-service catalog.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the input source.",
				Optional:            false,
				Computed:            false,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				}},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the input source",
				Optional:            true,
				Computed:            false,
				Required:            false,
			},
			"all_spaces": schema.BoolAttribute{
				Description: "Specify if is overridable at the blueprint level",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Validators: []validator.Bool{
					// Validate only this attribute or other_attr is configured or neither.
					boolvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("specific_spaces"),
					}...),
				},
			},
			"specific_spaces": schema.ListAttribute{
				Description: "Bucket's Name",
				Required:    false,
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					// Validate only this attribute or other_attr is configured or neither.
					listvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("all_spaces"),
					}...),
				},
			},
			"bucket_name_overridable": schema.BoolAttribute{
				Description: "Specify if is overridable at the blueprint level",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"bucket_name": schema.StringAttribute{
				Description: "Bucket's Name",
				Required:    true,
			},
			"filter_pattern_overridable": schema.BoolAttribute{
				Description: "Specify if is overridable at the blueprint level",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"filter_pattern": schema.StringAttribute{
				Description: "Bucket's Name",
				Required:    true,
			},
			"path_prefix_overridable": schema.BoolAttribute{
				Description: "Specify if is overridable at the blueprint level",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"path_prefix": schema.StringAttribute{
				Description: "Bucket's Name",
				Required:    true,
			},
			"credential_name": schema.StringAttribute{
				MarkdownDescription: "Credentials to use to connect to the bucket ",
				Optional:            false,
				Computed:            false,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// "all_spaces": schema.SingleNestedAttribute{
			// 	Description: "Environment collaborators",
			// 	Required:    false,
			// 	Optional:    true,
			// 	Attributes: map[string]schema.Attribute{
			// 		"all_spaces": schema.BoolAttribute{
			// 			Description: "Specify if is overridable at the blueprint level",
			// 			Optional:    true,
			// 			Computed:    true,
			// 			Default:     booldefault.StaticBool(true),
			// 		},
			// 		"specific_spaces": schema.ListAttribute{
			// 			Description: "Bucket's Name",
			// 			Required:    false,
			// 			Optional:    true,
			// 			ElementType: types.StringType,
			// 		},
			// 	},
			// },
			// "bucket_name": schema.SingleNestedAttribute{
			// 	Description: "Environment collaborators",
			// 	Required:    true,
			// 	Attributes: map[string]schema.Attribute{
			// 		"overridable": schema.BoolAttribute{
			// 			Description: "Specify if is overridable at the blueprint level",
			// 			Optional:    true,
			// 			Computed:    true,
			// 			Default:     booldefault.StaticBool(false),
			// 		},
			// 		"value": schema.StringAttribute{
			// 			Description: "Bucket's Name",
			// 			Required:    true,
			// 		},
			// 	},
			// },
			// "credential_name": schema.StringAttribute{
			// 	MarkdownDescription: "Credentials to use to connect to the bucket ",
			// 	Optional:            false,
			// 	Computed:            false,
			// 	Required:            true,
			// 	PlanModifiers: []planmodifier.String{
			// 		stringplanmodifier.RequiresReplace(),
			// 	},
			// },
			// "type": schema.StringAttribute{
			// 	MarkdownDescription: "Credentials to use to connect to the bucket ",
			// 	Optional:            false,
			// 	Computed:            false,
			// 	Required:            true,
			// 	PlanModifiers: []planmodifier.String{
			// 		stringplanmodifier.RequiresReplace(),
			// 	},
			// 	Validators: []validator.String{
			// 		stringvalidator.OneOf([]string{"s3-object", "s3-object-content"}...),
			// 	},
			// },
			// "filter_pattern": schema.SingleNestedAttribute{
			// 	Description: "Regex pattern to filter by",
			// 	Optional:    true,
			// 	Attributes: map[string]schema.Attribute{
			// 		"overridable": schema.BoolAttribute{
			// 			Description: "Specify if is overridable at the blueprint level",
			// 			Optional:    true,
			// 			Computed:    true,
			// 			Default:     booldefault.StaticBool(false),
			// 		},
			// 		"value": schema.StringAttribute{
			// 			Description: "Filter pattern value",
			// 			Required:    true,
			// 		},
			// 	},
			// },
			// "path_prefix": schema.SingleNestedAttribute{
			// 	Description: "Prefix of the path to get the input from.",
			// 	Optional:    true,
			// 	Attributes: map[string]schema.Attribute{
			// 		"overridable": schema.BoolAttribute{
			// 			Description: "Specify if is overridable at the blueprint level",
			// 			Optional:    true,
			// 			Computed:    true,
			// 			Default:     booldefault.StaticBool(false),
			// 		},
			// 		"value": schema.StringAttribute{
			// 			Description: "Path prefix value",
			// 			Required:    true,
			// 		},
			// 	},
			// },
		},
	}
}
func (r *TorqueS3ObjectInputSourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueS3ObjectInputSourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueS3ObjectInputSourceResourceModel
	var details client.InputSourceDetails
	var allowed_spaces client.AllowedSpaces
	const input_source_type = "s3-object"
	var specificSpaces []string
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	allowed_spaces.AllSpaces = data.AllSpaces.ValueBool()
	if !data.SpecificSpaces.IsNull() {
		allowed_spaces.AllSpaces = false
		for _, val := range data.SpecificSpaces.Elements() {
			specificSpaces = append(specificSpaces, strings.Replace(val.String(), "\"", "", -1))
		}
		allowed_spaces.SpecificSpaces = specificSpaces
	} else {
		allowed_spaces.AllSpaces = data.AllSpaces.ValueBool() // true
	}
	details.BucketName.Overridable = data.BucketNameOverridable.ValueBool()
	details.BucketName.Value = data.BucketName.ValueString()
	details.FilterPattern.Overridable = data.FilterPatternOverridable.ValueBool()
	details.FilterPattern.Value = data.FilterPattern.ValueString()
	details.PathPrefix.Overridable = data.PathPrefixOverridable.ValueBool()
	details.PathPrefix.Value = data.PathPrefix.ValueString()
	details.Type = input_source_type
	details.CredentialName = data.CredentialName.ValueString()
	err := r.client.CreateS3InputSource(data.Name.ValueString(), data.Description.ValueString(), allowed_spaces, details)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Input Source, got error: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueS3ObjectInputSourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

}

func (r *TorqueS3ObjectInputSourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueS3ObjectInputSourceResourceModel
	var state TorqueS3ObjectInputSourceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// if data.SelfService.ValueBool() {
	// 	err := r.client.PublishBlueprintInSpace(data.SpaceName.ValueString(), data.RepoName.ValueString(), data.Name.ValueString())
	// 	if err != nil {
	// 		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to publish workflow to self-service catalog, got error: %s", err))
	// 		return
	// 	}
	// } else {
	// 	err := r.client.UnpublishBlueprintInSpace(data.SpaceName.ValueString(), data.RepoName.ValueString(), data.Name.ValueString())
	// 	if err != nil {
	// 		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to unpublish workflow to self-service catalog, got error: %s", err))
	// 		return
	// 	}
	// }
	// if data.CustomIcon.IsNull() {
	// 	err := r.client.SetCatalogItemIcon(data.SpaceName.ValueString(), data.Name.ValueString(), data.RepoName.ValueString(), default_icon)
	// 	if err != nil {
	// 		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to remove Catalog Item custom icon, failed to set catalog item custom icon, got error: %s", err))
	// 		return
	// 	}
	// } else {
	// 	if data.CustomIcon.ValueString() != state.CustomIcon.ValueString() {
	// 		err := r.client.SetCatalogItemCustomIcon(data.SpaceName.ValueString(), data.Name.ValueString(), data.RepoName.ValueString(), data.CustomIcon.ValueString())
	// 		if err != nil {
	// 			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update workflow custom icon, failed to set catalog item custom icon, got error: %s", err))
	// 			return
	// 		}
	// 	}
	// }
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueS3ObjectInputSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueS3ObjectInputSourceResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteInputSource(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Input Source, got error: %s", err))
		return
	}

}

func (r *TorqueS3ObjectInputSourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
