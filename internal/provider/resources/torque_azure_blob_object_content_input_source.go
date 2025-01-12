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
var _ resource.Resource = &TorqueAzureBlobObjectContentInputSourceResource{}
var _ resource.ResourceWithImportState = &TorqueAzureBlobObjectContentInputSourceResource{}

func NewTorqueAzureBlobObjectContentInputSourceResource() resource.Resource {
	return &TorqueAzureBlobObjectContentInputSourceResource{}
}

// TorqueAzureBlobObjectContentInputSourceResource defines the resource implementation.
type TorqueAzureBlobObjectContentInputSourceResource struct {
	client *client.Client
}

// TorqueAzureBlobObjectContentInputSourceResourceModel describes the resource data model.
type TorqueAzureBlobObjectContentInputSourceResourceModel struct {
	Name                          types.String `tfsdk:"name"`
	Description                   types.String `tfsdk:"description"`
	AllSpaces                     types.Bool   `tfsdk:"all_spaces"`
	SpecificSpaces                types.List   `tfsdk:"specific_spaces"`
	StorageAccountName            types.String `tfsdk:"storage_account_name"`
	StorageAccountNameOverridable types.Bool   `tfsdk:"storage_account_overridable"`
	ContainerName                 types.String `tfsdk:"container_name"`
	ContainerNameOverridable      types.Bool   `tfsdk:"container_name_overridable"`
	BlobName                      types.String `tfsdk:"blob_name"`
	BlobNameOverrdiable           types.Bool   `tfsdk:"blob_name_overridable"`
	CredentialName                types.String `tfsdk:"credential_name"`
	FilterPattern                 types.String `tfsdk:"filter_pattern"`
	FilterPatternOverridable      types.Bool   `tfsdk:"filter_pattern_overridable"`
	JsonPath                      types.String `tfsdk:"json_path"`
	JsonPathOverridable           types.Bool   `tfsdk:"json_path_overridable"`
	DisplayJsonPath               types.String `tfsdk:"display_json_path"`
	DisplayJsonPathOverridable    types.Bool   `tfsdk:"display_json_path_overridable"`
}

func (r *TorqueAzureBlobObjectContentInputSourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_azure_blob_object_content_input_source"
}

func (r *TorqueAzureBlobObjectContentInputSourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Creates a new AWS Input Source of type Azure Blob Object.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the input source.",
				Optional:            false,
				Computed:            false,
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the input source",
				Optional:            true,
				Computed:            false,
				Required:            false,
			},
			"all_spaces": schema.BoolAttribute{
				Description: "Specify if the input source can be used in all spaces. Defaults to true, use specific spaces attribute for allowing only specific spaces.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					&allSpacesModifier{},
				},
				// Default:     booldefault.StaticBool(true),
				Validators: []validator.Bool{
					// Validate only this attribute or other_attr is configured or neither.
					boolvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("specific_spaces"),
					}...),
				},
			},
			"specific_spaces": schema.ListAttribute{
				Description: "List of spaces that can use this input source",
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
			"storage_account_overridable": schema.BoolAttribute{
				Description: "Specify if is overridable at the blueprint level",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"storage_account_name": schema.StringAttribute{
				Description: "Bucket's Name",
				Required:    true,
			},
			"container_name_overridable": schema.BoolAttribute{
				Description: "Specify if is overridable at the blueprint level",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"container_name": schema.StringAttribute{
				Description: "Bucket's Name",
				Required:    true,
			},
			"blob_name_overridable": schema.BoolAttribute{
				Description: "Specify if is overridable at the blueprint level",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"blob_name": schema.StringAttribute{
				Description: "Bucket's Name",
				Required:    true,
			},
			"filter_pattern_overridable": schema.BoolAttribute{
				Description: "Specify if is overridable at the blueprint level",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"filter_pattern": schema.StringAttribute{
				Description: "Regex pattern to filter the results by.",
				Required:    false,
				Optional:    true,
			},
			"json_path_overridable": schema.BoolAttribute{
				Description: "Specify if is overridable at the blueprint level",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"json_path": schema.StringAttribute{
				Description: "Enter the JSONPath for extracting the desired values",
				Required:    true,
				Optional:    false,
			},
			"display_json_path": schema.StringAttribute{
				Description: "Enter the JSONPath for extracting the corresponding display values",
				Required:    false,
				Optional:    true,
			},
			"display_json_path_overridable": schema.BoolAttribute{
				Description: "Specify if is overridable at the blueprint level.",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"credential_name": schema.StringAttribute{
				MarkdownDescription: "Credentials to use to connect to the bucket. Must be of type AWS.",
				Optional:            false,
				Computed:            false,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}
func (r *TorqueAzureBlobObjectContentInputSourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueAzureBlobObjectContentInputSourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueAzureBlobObjectContentInputSourceResourceModel
	var details client.InputSourceDetails
	const content_type = "JSON"
	details.ContentFormat = &client.ContentFormat{}         // pointer initialization       // pointer initialization
	details.StorageAccountName = &client.OverridableValue{} // pointer initialization
	details.ContainerName = &client.OverridableValue{}      // pointer initialization
	details.BlobName = &client.OverridableValue{}
	var allowed_spaces client.AllowedSpaces
	const input_source_type = "azure-blob-content"
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
	details.StorageAccountName.Overridable = data.StorageAccountNameOverridable.ValueBool()
	details.StorageAccountName.Value = data.StorageAccountName.ValueString()
	details.ContainerName.Overridable = data.ContainerNameOverridable.ValueBool()
	details.ContainerName.Value = data.ContainerName.ValueString()
	details.BlobName.Overridable = data.BlobNameOverrdiable.ValueBool()
	details.BlobName.Value = data.BlobName.ValueString()
	details.FilterPattern.Overridable = data.FilterPatternOverridable.ValueBool()
	details.FilterPattern.Value = data.FilterPattern.ValueString()
	details.ContentFormat.DisplayJsonPath.Value = data.DisplayJsonPath.ValueString()
	details.ContentFormat.DisplayJsonPath.Overridable = data.DisplayJsonPathOverridable.ValueBool()
	details.ContentFormat.JsonPath.Overridable = data.JsonPathOverridable.ValueBool()
	details.ContentFormat.JsonPath.Value = data.JsonPath.ValueString()
	details.ContentFormat.Type = content_type

	details.Type = input_source_type
	details.CredentialName = data.CredentialName.ValueString()
	err := r.client.CreateInputSource(data.Name.ValueString(), data.Description.ValueString(), allowed_spaces, details)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Input Source, got error: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueAzureBlobObjectContentInputSourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueAzureBlobObjectContentInputSourceResourceModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	input_source, err := r.client.GetInputSource(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Input Source details",
			"Could not read Input Source "+data.Name.ValueString()+": "+err.Error(),
		)
		return
	}
	data.Name = types.StringValue(input_source.Name)
	data.Description = types.StringValue(input_source.Description)
	data.StorageAccountName = types.StringValue(input_source.Details.StorageAccountName.Value)
	data.StorageAccountNameOverridable = types.BoolValue(input_source.Details.StorageAccountName.Overridable)
	data.ContainerName = types.StringValue(input_source.Details.ContainerName.Value)
	data.ContainerNameOverridable = types.BoolValue(input_source.Details.ContainerName.Overridable)
	data.CredentialName = types.StringValue(input_source.Details.CredentialName)
	data.AllSpaces = types.BoolValue(input_source.AllowedSpaces.AllSpaces)
	if len(input_source.AllowedSpaces.SpecificSpaces) > 0 {
		data.SpecificSpaces, _ = types.ListValueFrom(ctx, types.StringType, input_source.AllowedSpaces.SpecificSpaces)
	} else {
		data.SpecificSpaces = types.ListNull(types.StringType)
	}
	data.FilterPattern = types.StringValue(input_source.Details.FilterPattern.Value)
	data.FilterPatternOverridable = types.BoolValue(input_source.Details.FilterPattern.Overridable)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TorqueAzureBlobObjectContentInputSourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueAzureBlobObjectContentInputSourceResourceModel
	var state TorqueAzureBlobObjectContentInputSourceResourceModel
	var details client.InputSourceDetails

	details.PathPrefix = &client.OverridableValue{}         // pointer initialization
	details.StorageAccountName = &client.OverridableValue{} // pointer initialization
	details.ContainerName = &client.OverridableValue{}      // pointer initialization

	var allowed_spaces client.AllowedSpaces
	const input_source_type = "azure-blob-content"
	var specificSpaces []string
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

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
	details.StorageAccountName.Overridable = data.StorageAccountNameOverridable.ValueBool()
	details.StorageAccountName.Value = data.StorageAccountName.ValueString()
	details.ContainerName.Overridable = data.ContainerNameOverridable.ValueBool()
	details.ContainerName.Value = data.ContainerName.ValueString()
	details.FilterPattern.Overridable = data.FilterPatternOverridable.ValueBool()
	details.FilterPattern.Value = data.FilterPattern.ValueString()
	details.Type = input_source_type
	details.CredentialName = data.CredentialName.ValueString()
	err := r.client.UpdateInputSource(state.Name.ValueString(), data.Name.ValueString(), data.Description.ValueString(), allowed_spaces, details)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Input Source, got error: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueAzureBlobObjectContentInputSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueAzureBlobObjectContentInputSourceResourceModel

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

func (r *TorqueAzureBlobObjectContentInputSourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// type allSpacesModifier struct{}

// func (m allSpacesModifier) Description(ctx context.Context) string {
// 	return "Set 'all_spaces' to false if 'specific_spaces' is provided and non-empty."
// }

// func (m allSpacesModifier) MarkdownDescription(ctx context.Context) string {
// 	return "Set `all_spaces` to `false` if `specific_spaces` is provided and non-empty."
// }

// func (m allSpacesModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
// 	// If the user explicitly set 'all_spaces', respect that value.
// 	if !req.ConfigValue.IsNull() {
// 		resp.PlanValue = req.ConfigValue
// 		return
// 	}

// 	// Retrieve 'specific_spaces' from the planned state.
// 	var specificSpaces []string
// 	diags := req.Plan.GetAttribute(ctx, path.Root("specific_spaces"), &specificSpaces)
// 	if diags.HasError() {
// 		resp.Diagnostics.Append(diags...)
// 		return
// 	}

// 	// If 'specific_spaces' is non-empty, set 'all_spaces' to false.
// 	if len(specificSpaces) > 0 {
// 		resp.PlanValue = types.BoolValue(false)
// 	} else {
// 		// Otherwise, default to true.
// 		resp.PlanValue = types.BoolValue(true)
// 	}
// }
