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
var _ resource.Resource = &TorqueAzureBlobObjectInputSourceResource{}
var _ resource.ResourceWithImportState = &TorqueAzureBlobObjectInputSourceResource{}

func NewTorqueAzureBlobObjectInputSourceResource() resource.Resource {
	return &TorqueAzureBlobObjectInputSourceResource{}
}

// TorqueAzureBlobObjectInputSourceResource defines the resource implementation.
type TorqueAzureBlobObjectInputSourceResource struct {
	client *client.Client
}

// TorqueAzureBlobObjectInputSourceResourceModel describes the resource data model.
type TorqueAzureBlobObjectInputSourceResourceModel struct {
	Name                          types.String `tfsdk:"name"`
	Description                   types.String `tfsdk:"description"`
	AllSpaces                     types.Bool   `tfsdk:"all_spaces"`
	SpecificSpaces                types.List   `tfsdk:"specific_spaces"`
	StorageAccountName            types.String `tfsdk:"storage_account_name"`
	StorageAccountNameOverridable types.Bool   `tfsdk:"storage_account_overridable"`
	ContainerName                 types.String `tfsdk:"container_name"`
	ContainerNameOverridable      types.Bool   `tfsdk:"container_name_overridable"`
	CredentialName                types.String `tfsdk:"credential_name"`
	FilterPattern                 types.String `tfsdk:"filter_pattern"`
	FilterPatternOverridable      types.Bool   `tfsdk:"filter_pattern_overridable"`
	PathPrefix                    types.String `tfsdk:"path_prefix"`
	PathPrefixOverridable         types.Bool   `tfsdk:"path_prefix_overridable"`
}

func (r *TorqueAzureBlobObjectInputSourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_azure_blob_object_input_source"
}

func (r *TorqueAzureBlobObjectInputSourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"path_prefix_overridable": schema.BoolAttribute{
				Description: "Specify if is overridable at the blueprint level",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"path_prefix": schema.StringAttribute{
				Description: "Path prefix of the object.",
				Required:    false,
				Optional:    true,
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
func (r *TorqueAzureBlobObjectInputSourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueAzureBlobObjectInputSourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueAzureBlobObjectInputSourceResourceModel
	var details client.InputSourceDetails

	details.PathPrefix = &client.OverridableValue{}         // pointer initialization
	details.StorageAccountName = &client.OverridableValue{} // pointer initialization
	details.ContainerName = &client.OverridableValue{}      // pointer initialization

	var allowed_spaces client.AllowedSpaces
	const input_source_type = "azure-blob"
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
	details.FilterPattern.Overridable = data.FilterPatternOverridable.ValueBool()
	details.FilterPattern.Value = data.FilterPattern.ValueString()
	details.PathPrefix.Overridable = data.PathPrefixOverridable.ValueBool()
	details.PathPrefix.Value = data.PathPrefix.ValueString()
	details.Type = input_source_type
	details.CredentialName = data.CredentialName.ValueString()
	err := r.client.CreateInputSource(data.Name.ValueString(), data.Description.ValueString(), allowed_spaces, details)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Input Source, got error: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueAzureBlobObjectInputSourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueAzureBlobObjectInputSourceResourceModel

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
	data.PathPrefix = types.StringValue(input_source.Details.PathPrefix.Value)
	data.PathPrefixOverridable = types.BoolValue(input_source.Details.PathPrefix.Overridable)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TorqueAzureBlobObjectInputSourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueAzureBlobObjectInputSourceResourceModel
	var state TorqueAzureBlobObjectInputSourceResourceModel
	var details client.InputSourceDetails

	details.PathPrefix = &client.OverridableValue{}         // pointer initialization
	details.StorageAccountName = &client.OverridableValue{} // pointer initialization
	details.ContainerName = &client.OverridableValue{}      // pointer initialization

	var allowed_spaces client.AllowedSpaces
	const input_source_type = "azure-blob"
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
	details.PathPrefix.Overridable = data.PathPrefixOverridable.ValueBool()
	details.PathPrefix.Value = data.PathPrefix.ValueString()
	details.Type = input_source_type
	details.CredentialName = data.CredentialName.ValueString()
	err := r.client.UpdateInputSource(state.Name.ValueString(), data.Name.ValueString(), data.Description.ValueString(), allowed_spaces, details)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Input Source, got error: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueAzureBlobObjectInputSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueAzureBlobObjectInputSourceResourceModel

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

func (r *TorqueAzureBlobObjectInputSourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
