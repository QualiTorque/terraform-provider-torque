package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/qualitorque/terraform-provider-torque/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TorqueBlueprintResource{}
var _ resource.ResourceWithImportState = &TorqueBlueprintResource{}

func NewTorqueBlueprintResource() resource.Resource {
	return &TorqueBlueprintResource{}
}

// TorqueBlueprintResource defines the resource implementation.
type TorqueBlueprintResource struct {
	client *client.Client
}

// TorqueBlueprintResourceModel describes the resource data model.
type TorqueBlueprintResourceModel struct {
	BlueprintName types.String `tfsdk:"blueprint_name"`
	Name          types.String `tfsdk:"name"`
	DisplayName   types.String `tfsdk:"display_name"`
	RepoName      types.String `tfsdk:"repository_name"`
	RepoBranch    types.String `tfsdk:"repository_branch"`
	Commit        types.String `tfsdk:"commit"`
	Description   types.String `tfsdk:"description"`
	Url           types.String `tfsdk:"url"`
	ModifiedBy    types.String `tfsdk:"modified_by"`
	Published     types.Bool   `tfsdk:"enabled"`
}

func (r *TorqueBlueprintResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_blueprint"
}

func (r *TorqueBlueprintResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Creation of a new Torque space with associated entities (users, repos, etc...)",

		Attributes: map[string]schema.Attribute{
			"blueprint_name": schema.StringAttribute{
				Description: "The unique name of the blueprint",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "no idea!! ",
				Optional:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "The user friendly name of the blueprint in the space",
				Computed:    true,
			},
			"repository_name": schema.StringAttribute{
				Description: "The repository name from which the blueprint is used",
				Computed:    true,
			},
			"repository_branch": schema.StringAttribute{
				Description: "The branch from which the blueprint is taken",
				Computed:    true,
			},
			"commit": schema.StringAttribute{
				Description: "The commit id of the blueprint",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the bluepring",
				Computed:    true,
			},
			"url": schema.StringAttribute{
				Description: "URI of the blueprint",
				Computed:    true,
			},
			"modified_by": schema.StringAttribute{
				Description: "The name of the user that last modified the blueprint",
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "is Published blueprint in the space",
				Computed:    true,
			},
		},
	}
}

func (r *TorqueBlueprintResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueBlueprintResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueBlueprintResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueBlueprintResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueBlueprintResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	results, err := r.client.GetSpaceBlueprints(r.client.Space)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read blueprint , got error: %s", err))
		return
	}

	for _, n := range results {
		if data.Name.ValueString() == n.BlueprintName {
			data.BlueprintName = types.StringValue(n.BlueprintName)
			//data.Name = types.StringValue(n.Name)
			data.RepoName = types.StringValue(n.RepoName)
			data.Description = types.StringValue(n.Description)
			data.Commit = types.StringValue(n.Commit)
			data.ModifiedBy = types.StringValue(n.ModifiedBy)
			data.DisplayName = types.StringValue(n.DisplayName)
			data.RepoBranch = types.StringValue(n.RepoBranch)
			data.Url = types.StringValue(n.Url)
			data.Published = types.BoolValue(n.Published)
			break
		}
	}

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TorqueBlueprintResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueBlueprintResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueBlueprintResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueBlueprintResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TorqueBlueprintResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
