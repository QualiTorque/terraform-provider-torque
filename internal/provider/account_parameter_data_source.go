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
	_ datasource.DataSource              = &accountParameterDataSource{}
	_ datasource.DataSourceWithConfigure = &accountParameterDataSource{}
)

// NewaccountParametersDataSource is a helper function to simplify the provider implementation.
func NewAccountParameterDataSource() datasource.DataSource {
	return &accountParameterDataSource{}
}

// accountParameterDataSource is the data source implementation.
type accountParameterDataSource struct {
	client *client.Client
}

// accountParameterDataSourceModel maps the data source schema data.
type accountParameterDataSourceModel struct {
	Name        types.String `tfsdk:"name"`
	Value       types.String `tfsdk:"value"`
	Sensitive   types.Bool   `tfsdk:"sensitive"`
	Description types.String `tfsdk:"description"`
}

// Metadata returns the data source type name.
func (d *accountParameterDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account_parameter"
}

// Schema defines the schema for the data source.
func (d *accountParameterDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get Account Parameter details.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the account level parameter",
				Required:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "Parameter Value. Value will be an empty string if the parameter is sensitive",
				Computed:            true,
			},
			"sensitive": schema.BoolAttribute{
				MarkdownDescription: "Whether the parameter is sensitive or not.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Parameter description",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *accountParameterDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *accountParameterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state accountParameterDataSourceModel

	var parameter types.String

	diags1 := req.Config.GetAttribute(ctx, path.Root("name"), &parameter)
	resp.Diagnostics.Append(diags1...)
	if resp.Diagnostics.HasError() {
		return
	}

	account_parameter_data, err := d.client.GetAccountParameter(parameter.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Torque Account Parameter",
			err.Error(),
		)
		return
	}

	state.Name = types.StringValue(account_parameter_data.Name)
	state.Value = types.StringValue(account_parameter_data.Value)
	state.Sensitive = types.BoolValue(account_parameter_data.Sensitive)
	state.Description = types.StringValue(account_parameter_data.Description)
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
