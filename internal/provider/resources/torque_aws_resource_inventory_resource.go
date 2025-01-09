package resources

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/qualitorque/terraform-provider-torque/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TorqueAwsResourceInventoryResource{}
var _ resource.ResourceWithImportState = &TorqueAwsResourceInventoryResource{}

func NewTorqueAwsResourceInventoryResource() resource.Resource {
	return &TorqueAwsResourceInventoryResource{}
}

// TorqueAwsResourceInventoryResource defines the resource implementation.
type TorqueAwsResourceInventoryResource struct {
	client *client.Client
}

// TorqueAwsResourceInventoryResourceModel describes the resource data model.
type TorqueAwsResourceInventoryResourceModel struct {
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	AccountNumber types.String `tfsdk:"account_number"`
	AccessKey     types.String `tfsdk:"access_key"`
	SecretKey     types.String `tfsdk:"secret_key"`
	// Credentials types.String `tfsdk:"credentials"`
	ViewArn   types.String `tfsdk:"view_arn"`
	CloudType types.String `tfsdk:"cloud_type"`
}

func (r *TorqueAwsResourceInventoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_aws_resource_inventory"
}

func (r *TorqueAwsResourceInventoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Creation of a new parameter is a Torque space",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the AWS Cloud Account. Will also be used to store the provided credentials in Torque's credential store for later use.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "ARN of the reader IAM role.",
				Required:            true,
			},
			"account_number": schema.StringAttribute{
				MarkdownDescription: "ARN of the reader IAM role.",
				Required:            true,
			},
			"access_key": schema.StringAttribute{
				MarkdownDescription: "AWS Access Key.",
				Required:            true,
				Sensitive:           true,
			},
			"secret_key": schema.StringAttribute{
				MarkdownDescription: "AWS Secret Access Key.",
				Required:            true,
				Sensitive:           true,
			},
			"view_arn": schema.StringAttribute{
				MarkdownDescription: "ARN of the reader IAM role.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^arn:(aws|aws-cn|aws-us-gov):[a-zA-Z0-9\-]+:[a-zA-Z0-9\-]*:[0-9]{12}:[a-zA-Z0-9\-_/\.]+$`),
						"must be a valid ARN format (e.g., arn:aws:iam::123456789012:user/JohnDoe)",
					),
				},
			},
			"cloud_type": schema.StringAttribute{
				MarkdownDescription: "Type of the resource inventory.",
				Required:            false,
				Computed:            true,
				Default:             stringdefault.StaticString("aws"),
			},
		},
	}
}

func (r *TorqueAwsResourceInventoryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueAwsResourceInventoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueAwsResourceInventoryResourceModel
	var details client.ResourceInventoryDetails
	const credential_type = "aws__basic"
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.CreateAccountCredentials(data.Name.ValueString(), data.Description.ValueString(), data.CloudType.ValueString(), data.AccountNumber.ValueString(), credential_type, nil, data.AccessKey.ValueStringPointer(), data.SecretKey.ValueStringPointer(), nil)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create resource inventory credentials, got error: %s", err))
		return
	}
	details.Type = data.CloudType.ValueString()
	details.ViewArn = data.ViewArn.ValueStringPointer()
	err = r.client.ConfigureResourveInventory(data.Name.ValueString(), details)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create resource inventory, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Resource Created Successful!")

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueAwsResourceInventoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueAwsResourceInventoryResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueAwsResourceInventoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueAwsResourceInventoryResourceModel
	const credential_type = "aws__basic"
	var details client.ResourceInventoryDetails

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	err := r.client.UpdateAccountCredentials(data.Name.ValueString(), data.Description.ValueString(), data.AccountNumber.ValueString(), data.CloudType.ValueString(), credential_type, nil, data.AccessKey.ValueStringPointer(), data.SecretKey.ValueStringPointer(), nil)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update resource inventory credentials, got error: %s", err))
		return
	}
	details.Type = data.CloudType.ValueString()
	details.ViewArn = data.ViewArn.ValueStringPointer()

	err = r.client.ConfigureResourveInventory(data.Name.ValueString(), details)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update resource inventory, got error: %s", err))
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueAwsResourceInventoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueAwsResourceInventoryResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	err := r.client.DeleteResourceInventory(data.Name.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create resource inventory, got error: %s", err))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *TorqueAwsResourceInventoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("credentials"), req, resp)
}
