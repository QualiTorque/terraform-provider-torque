package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	Yaml                types.String `tfsdk:"yaml"`
	DisplayName         types.String `tfsdk:"display_name"`
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
		Description: "Get details of an account level workflow.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the workflow",
				Required:            true,
			},
			"yaml": schema.StringAttribute{
				MarkdownDescription: "Yaml formatted string that describes the workflow.",
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the workflow.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Workflow description",
				Computed:            true,
			},
			"enforced_on_all_spaces": schema.BoolAttribute{
				MarkdownDescription: "Whether the workflow is enforced on all spaces or not.",
				Computed:            true,
			},
			"specific_spaces": schema.ListAttribute{
				MarkdownDescription: "List of spaces the workflow is enforced on if enforced on all spaces is false. Empty list if enforced_on_all_spaces is true",
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

	diags := req.Config.GetAttribute(ctx, path.Root("name"), &workflow_name)
	resp.Diagnostics.Append(diags...)
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
	state.Name = types.StringValue(workflow_data.Name)
	state.Yaml = types.StringValue(workflow_data.Yaml)
	state.DisplayName = types.StringValue(workflow_data.DisplayName)
	state.Description = types.StringValue(workflow_data.Description)
	state.EnforcedOnAllSpaces = types.BoolValue(workflow_data.SpaceDefinition.EnforcedOnAllSpaces)
	if workflow_data.SpaceDefinition.EnforcedOnAllSpaces {
		state.SpecificSpaces = types.ListValueMust(types.StringType, []attr.Value{})
	} else {
		state.SpecificSpaces, _ = types.ListValueFrom(ctx, types.StringType, workflow_data.SpaceDefinition.SpecificSpaces)
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
