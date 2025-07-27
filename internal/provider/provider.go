package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/qualitorque/terraform-provider-torque/client"
	"github.com/qualitorque/terraform-provider-torque/internal/provider/data_sources"
	"github.com/qualitorque/terraform-provider-torque/internal/provider/resources"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &torqueProvider{}

// TorqueProvider defines the provider implementation.
type torqueProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// torqueProviderModel describes the provider data model.
type torqueProviderModel struct {
	Host  types.String `tfsdk:"host"`
	Space types.String `tfsdk:"space"`
	Token types.String `tfsdk:"token"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &torqueProvider{
			version: version,
		}
	}
}

func (p *torqueProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "torque"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *torqueProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `
		Use Torque provider to interact with Torque by Quali.<br/><br/>
		Torque by Quali is an Environment-as-a-Service (EaaS) control plane and self-service catalog allowing you to deploy and manage cloud environments  
		comprising the infrastructure, applications,and any dependencies or external services necessary for applications or services to rely on.<br/><br/>
		For more information, visit [Quali's Torque Documentation](https://docs.qtorque.io)  
		For experimenting with Torque, visit [Quali's Torque Playground](https://www.quali.com/watch-see-the-torque-playground-in-action/)
		`,
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "URI for Torque API. May also be provided via TORQUE_HOST environment variable.",
				Optional:    true,
			},
			"space": schema.StringAttribute{
				Description: "Space for Torque API. May also be provided via TORQUE_SPACE environment variable.",
				Optional:    true,
			},
			"token": schema.StringAttribute{
				Description: "Token for Torque API. May also be provided via TORQUE_TOKEN environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares a HashiCups API client for data sources and resources.
func (p *torqueProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring torque-provider client")
	var config torqueProviderModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	// if config.Host.IsUnknown() {
	// 	resp.Diagnostics.AddAttributeError(
	// 		path.Root("host"),
	// 		"Unknown Torque API Host",
	// 		"The provider cannot create the Torque API client as there is an unknown configuration value for the Torque API host. "+
	// 			"Either target apply the source of the value first, set the value statically in the configuration, or use the TORQUE_HOST environment variable.",
	// 	)
	// }.

	// if config.Space.IsUnknown() {
	// 	resp.Diagnostics.AddAttributeError(
	// 		path.Root("space"),
	// 		"Unknown Torque API space",
	// 		"The provider cannot create the Torque API client as there is an unknown configuration value for the Torque API space. "+
	// 			"Either target apply the source of the value first, set the value statically in the configuration, or use the TORQUE_SPACE environment variable.",
	// 	)
	// }.

	// if config.Token.IsUnknown() {
	// 	resp.Diagnostics.AddAttributeError(
	// 		path.Root("token"),
	// 		"Unknown Torque API token",
	// 		"The provider cannot create the Torque API client as there is an unknown configuration value for the Torque API token or long-token. "+
	// 			"Either target apply the source of the value first, set the value statically in the configuration, or use the TORQUE_TOKEN environment variable.",
	// 	)
	// }.

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	host := os.Getenv("TORQUE_HOST")
	space := os.Getenv("TORQUE_SPACE")
	token := os.Getenv("TORQUE_TOKEN")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Space.IsNull() {
		space = config.Space.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	// if host == "" {
	// 	resp.Diagnostics.AddAttributeError(
	// 		path.Root("host"),
	// 		"Missing Torque API Host",
	// 		"The provider cannot create the Torque API client as there is a missing or empty value for the Torque API host. "+
	// 			"Set the host value in the configuration or use the TORQUE_HOST environment variable. "+
	// 			"If either is already set, ensure the value is not empty.",
	// 	)
	// }.

	// if space == "" {
	// 	resp.Diagnostics.AddAttributeError(
	// 		path.Root("space"),
	// 		"Missing Torque API space",
	// 		"The provider cannot create the Torque API client as there is a missing or empty value for the Torque API username. "+
	// 			"Set the username value in the configuration or use the TORQUE_SPACE environment variable. "+
	// 			"If either is already set, ensure the value is not empty.",
	// 	)
	// }.

	// if token == "" {
	// 	resp.Diagnostics.AddAttributeError(
	// 		path.Root("token"),
	// 		"Missing Torque API token",
	// 		"The provider cannot create the Torque API client as there is a missing or empty value for the Torque API password. "+
	// 			"Set the password value in the configuration or use the TORQUE_TOKEN environment variable. "+
	// 			"If either is already set, ensure the value is not empty.",
	// 	)
	// }.

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "torque_host", host)
	ctx = tflog.SetField(ctx, "torque_space", space)
	ctx = tflog.SetField(ctx, "torque_token", token)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "torque_token")

	tflog.Debug(ctx, "Creating Torque API client")

	// Create a new Torque client using the configuration values.
	client, err := client.NewClient(&host, &space, &token)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Torque API Client",
			"An unexpected error occurred when creating the Torque API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Torque Client Error: "+err.Error(),
		)
		return
	}

	// Make the Torque client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Torque client", map[string]any{"success": true})
}

func (p *torqueProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewTorqueIntrospectionResource,
		resources.NewTorqueAgentSpaceAssociationResource,
		resources.NewTorqueSpaceRepositoryResource,
		resources.NewTorqueSpaceResource,
		resources.NewTorqueUserSpaceAssociationResource,
		resources.NewTorqueTagResource,
		resources.NewTorqueCatalogItemResource,
		resources.NewTorqueTagSpaceValueAssociationResource,
		resources.NewTorqueParameterResource,
		resources.NewTorqueSpaceParameterResource,
		resources.NewTorqueGroupResource,
		resources.NewTorqueAwsCostTargetResource,
		resources.NewTorqueTagBlueprintValueAssociationResource,
		resources.NewTorqueSpaceEmailNotificationResource,
		resources.NewTorqueAccountResource,
		resources.NewTorqueSpaceCodeCommitRepositoryResource,
		resources.NewTorqueSpaceGitlabEnterpriseRepositoryResource,
		resources.NewTorqueAssetLibraryItemResource,
		resources.NewTorqueSpaceLabelResource,
		resources.NewTorqueSpaceLabelAssociationResource,
		resources.NewTorqueEnvironmentLabelResource,
		resources.NewTorqueEnvironmentLabelAssociationResource,
		resources.NewTorqueSpaceGitCredentialsResource,
		resources.NewTorqueGitCredentialsResource,
		resources.NewTorqueEnvironmentResource,
		resources.NewTorqueWorkflowResource,
		resources.NewTorqueSpaceWorkflowResource,
		resources.NewTorqueSpaceCustomIconResource,
		resources.NewTorqueS3ObjectInputSourceResource,
		resources.NewTorqueS3ObjectContentInputSourceResource,
		resources.NewTorqueAwsResourceInventoryResource,
		resources.NewTorqueAzureBlobObjectInputSourceResource,
		resources.NewTorqueAzureBlobObjectContentInputSourceResource,
		resources.NewTorqueDeploymentEngineResource,
		resources.NewTorqueEmailApprovalChannelResource,
		resources.NewTorqueTeamsApprovalChannelResource,
		resources.NewTorqueServiceNowApprovalChannelResource,
		resources.NewTorqueSpaceTeamsNotificationResource,
		resources.NewTorqueSpaceSlackNotificationResource,
		resources.NewTorqueSpaceGenericWebhookNotificationResource,
		resources.NewTorqueAuditResource,
		resources.NewTorqueElasticsearchAuditResource,
		resources.NewTorqueSpaceAdoServerRepositoryResource,
	}
}

// DataSources defines the data sources implemented in the provider.
func (p *torqueProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		data_sources.NewUserDataSource,
		data_sources.NewSpaceRepositoryBlueprintsDataSource,
		data_sources.NewEnvironmentDataSource,
		data_sources.NewEnvironmentIntrospectionDataSource,
		data_sources.NewAccountParameterDataSource,
		data_sources.NewSpaceParameterDataSource,
		data_sources.NewSpaceBlueprintDataSource,
		data_sources.NewTorqueWorkflowDataSource,
		data_sources.NewSpaceCustomIconDataSource,
		data_sources.NewSpacesDataSource,
	}
}
