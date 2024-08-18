package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/qualitorque/terraform-provider-torque/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &torqueWorkflow{}
	_ datasource.DataSourceWithConfigure = &torqueWorkflow{}
)

// NewaccountParametersDataSource is a helper function to simplify the provider implementation.
func NewTorqueWorkflowDataSource() datasource.DataSource {
	return &torqueWorkflow{}
}

// torqueWorkflow is the data source implementation.
type torqueWorkflow struct {
	client *client.Client
}

// torqueWorkflowModel maps the data source schema data.
type torqueWorkflowModel struct {
	Name                types.String `tfsdk:"name"`
	Yaml                types.String `tfsdk:"value"`
	DisplayName         types.String `tfsdk:"sensitive"`
	Description         types.String `tfsdk:"description"`
	EnforcedOnAllSpaces types.Bool   `tfsdk:"enforced_on_all_spaces"`
	SpecificSpaces      types.List   `tfsdk:"specific_spaces"`
}

// Metadata returns the data source type name.
func (d *torqueWorkflow) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow"
}

// Schema defines the schema for the data source.
func (d *torqueWorkflow) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get Account Parameter details.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the account level parameter",
				Required:            true,
			},
			"yaml": schema.StringAttribute{
				MarkdownDescription: "Parameter Value. Value of sensitive parameter is null.",
				Computed:            true,
			},
			"display_name": schema.BoolAttribute{
				MarkdownDescription: "Whether the parameter is sensitive or not.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Parameter description",
				Computed:            true,
			},
			"enforced_on_all_spaces": schema.BoolAttribute{
				MarkdownDescription: "Whether the parameter is sensitive or not.",
				Computed:            true,
			},
			"specific_spaces": schema.ListAttribute{
				MarkdownDescription: "Parameter description",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *torqueWorkflow) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *torqueWorkflow) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state torqueWorkflowModel

	var workflow_name types.String

	diags1 := req.Config.GetAttribute(ctx, path.Root("workflow_name"), &workflow_name)
	resp.Diagnostics.Append(diags1...)
	if resp.Diagnostics.HasError() {
		return
	}

	workflow_data, err := d.client.GetWorkflow(workflow_name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Torque Workflow",
			err.Error(),
		)
		return
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
