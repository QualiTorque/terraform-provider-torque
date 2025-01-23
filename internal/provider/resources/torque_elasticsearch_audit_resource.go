package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/qualitorque/terraform-provider-torque/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TorqueElasticsearchAuditResource{}
var _ resource.ResourceWithImportState = &TorqueElasticsearchAuditResource{}

func NewTorqueElasticsearchAuditResource() resource.Resource {
	return &TorqueElasticsearchAuditResource{}
}

// TorqueElasticsearchAuditResource defines the resource implementation.
type TorqueElasticsearchAuditResource struct {
	client *client.Client
}

// TorqueElasticsearchAuditResourceModel describes the resource data model.
type TorqueElasticsearchAuditResourceModel struct {
	Url         types.String `tfsdk:"url"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
	Certificate types.String `tfsdk:"certificate"`
	Type        types.String `tfsdk:"type"`
}

func (r *TorqueElasticsearchAuditResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_elasticsearch_audit"
}

func (r *TorqueElasticsearchAuditResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Creates an Elasticsearch Audit target. Once integrated, Torque begins capturing events and you’ll ship them to the configured Elasticsearch instance.",

		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				MarkdownDescription: "Elasticsearch instance URL.",
				Optional:            false,
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Elasticsearch instance username.",
				Optional:            false,
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Elasticsearch instance password.",
				Optional:            false,
				Sensitive:           true,
				Required:            true,
			},
			"certificate": schema.StringAttribute{
				MarkdownDescription: "Optional certificate of the Elasticsearch instance.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				Default:             stringdefault.StaticString(""),
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Built-in audit is of type elasticsearch.",
				Optional:            false,
				Computed:            true,
				Default:             stringdefault.StaticString("elasticsearch"),
			},
		},
	}
}

func (r *TorqueElasticsearchAuditResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueElasticsearchAuditResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueElasticsearchAuditResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	var properties client.AuditProperties
	properties.Username = data.Username.ValueString()
	properties.Password = data.Username.ValueString()
	properties.Url = data.Url.ValueString()
	properties.Certificate = data.Certificate.ValueStringPointer()
	err := r.client.CreateAuditTarget(data.Type.ValueString(), &properties)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to configure audit, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueElasticsearchAuditResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueElasticsearchAuditResourceModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	audit, err := r.client.GetAudit()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Audit Target configuration",
			"Could not read Audit Target: "+err.Error(),
		)
		return
	}

	data.Url = types.StringValue(audit.Properties.Url)
	data.Username = types.StringValue(audit.Properties.Username)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TorqueElasticsearchAuditResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueElasticsearchAuditResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Type = types.StringValue("elasticsearch")
	var properties client.AuditProperties
	properties.Username = data.Username.ValueString()
	properties.Password = data.Password.ValueString()
	properties.Url = data.Url.ValueString()
	properties.Certificate = data.Certificate.ValueStringPointer()
	err := r.client.CreateAuditTarget(data.Type.ValueString(), &properties)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to configure audit, got error: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueElasticsearchAuditResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueElasticsearchAuditResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAudit(data.Type.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete audit, got error: %s", err))
		return
	}

}

func (r *TorqueElasticsearchAuditResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
