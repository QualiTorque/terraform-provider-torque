package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/qualitorque/terraform-provider-torque/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TorqueAccountResource{}
var _ resource.ResourceWithImportState = &TorqueAccountResource{}

func NewTorqueAccountResource() resource.Resource {
	return &TorqueAccountResource{}
}

// TorqueAccountResource defines the resource implementation.
type TorqueAccountResource struct {
	client *client.Client
}

// POST  https://portal.qtorque.io/api/accounts/<PARENT_ACCOUNT_NAME>/subaccounts
// {"parent_account":"tomera","account_name":"account","password":"somepassword"}

// TorqueAccountResourceModel describes the resource data model.
type TorqueAccountResourceModel struct {
	ParentAccount types.String `tfsdk:"parent_account"`
	AccountName   types.String `tfsdk:"account_name"`
	Password      types.String `tfsdk:"password"`
	Company       types.String `tfsdk:"company"`
}

func (r *TorqueAccountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_account"
}

func (r *TorqueAccountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Creation of a new Torque account with associated entities (users, repos, etc...)",

		Attributes: map[string]schema.Attribute{
			"parent_account": schema.StringAttribute{
				MarkdownDescription: "Name of the new Account to be added to torque",
				Required:            true,
			},
			"account_name": schema.StringAttribute{
				MarkdownDescription: "Name of the new Account to be added to torque",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Account value to be set as the Account value default",
				Required:            true,
				Computed:            false,
			},
			"company": schema.StringAttribute{
				MarkdownDescription: "Account scope. Possible values: account, account, blueprint, environment",
				Required:            true,
				Computed:            false,
			},
		},
	}
}

func (r *TorqueAccountResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TorqueAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueAccountResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// var possible []string
	// if !data.PossibleValues.IsNull() {
	// 	for _, pos_value := range data.PossibleValues.Elements() {
	// 		possible = append(possible, pos_value.String())
	// 	}
	// }

	err := r.client.CreateAccount(data.ParentAccount.ValueString(), data.AccountName.ValueString(), data.Password.ValueString(), data.Company.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Account, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Resource Created Successful!")

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *TorqueAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueAccountResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueAccountResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueAccountResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the account.
	err := r.client.RemoveAccount(data.AccountName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Account, got error: %s", err))
		return
	}

}

func (r *TorqueAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
