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
var _ resource.Resource = &TorqueSpaceResource{}
var _ resource.ResourceWithImportState = &TorqueSpaceResource{}

func NewTorqueSpaceResource() resource.Resource {
	return &TorqueSpaceResource{}
}

// TorqueSpaceResource defines the resource implementation.
type TorqueSpaceResource struct {
	client *client.Client
}

// TorqueSpaceResourceModel describes the resource data model.
type TorqueSpaceResourceModel struct {
	Name              types.String `tfsdk:"name"`
	Color             types.String `tfsdk:"color"`
	Icon              types.String `tfsdk:"icon"`
	AssociatedMembers types.List   `tfsdk:"space_members"`
	AssociatedAdmins  types.List   `tfsdk:"space_admins"`
	AssociatedAgents  types.Map    `tfsdk:"associated_kubernetes_agent"`
}

func (r *TorqueSpaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "torque_space_resource"
}

func (r *TorqueSpaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Resource that will be presented in Torque resource catalog",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Space name to be create",
				Required:            true,
			},
			"color": schema.StringAttribute{
				MarkdownDescription: "Space color to be used for the new space",
				Optional:            true,
				Computed:            false,
			},
			"icon": schema.StringAttribute{
				MarkdownDescription: "Space icon to be used",
				Optional:            true,
				Computed:            false,
			},
			"space_members": schema.ListAttribute{
				MarkdownDescription: "List of space members to be associate to the newly created space",
				Optional:            true,
				Computed:            false,
				ElementType:         types.StringType,
			},
			"space_admins": schema.ListAttribute{
				MarkdownDescription: "List of space admins to be associate to the newly created space",
				Optional:            true,
				Computed:            false,
				ElementType:         types.StringType,
			},
			"associated_kubernetes_agent": schema.MapAttribute{
				MarkdownDescription: "Kubernetes agent to associate to the newly create space",
				Optional:            true,
				Computed:            false,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *TorqueSpaceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func trimQuote(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}

func (r *TorqueSpaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TorqueSpaceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var name types.String
	req.Config.GetAttribute(ctx, path.Root("name"), &name)

	var color types.String
	req.Config.GetAttribute(ctx, path.Root("color"), &color)

	var icon types.String
	req.Config.GetAttribute(ctx, path.Root("icon"), &icon)

	var space_memebrs types.List
	req.Config.GetAttribute(ctx, path.Root("space_members"), &space_memebrs)

	var space_admins types.List
	req.Config.GetAttribute(ctx, path.Root("space_admins"), &space_admins)

	var associated_agents types.List
	req.Config.GetAttribute(ctx, path.Root("associated_agents"), &associated_agents)

	var agent types.Map
	req.Config.GetAttribute(ctx, path.Root("associated_kubernetes_agent"), &agent)

	err := r.client.CreateSpace(name.ValueString(), color.ValueString(), icon.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create space, got error: %s", err))
		return
	}

	data.Name = name
	data.Color = color
	data.Icon = icon

	if space_memebrs.IsNull() == false {
		for _, member := range space_memebrs.Elements() {
			err := r.client.AddUserToSpace(trimQuote(member.String()), "Space Member", name.ValueString())
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to attach space member to space, got error: %s", err))
				return
			}
		}
	}

	if space_admins.IsNull() == false {
		for _, admin := range space_admins.Elements() {
			err := r.client.AddUserToSpace(trimQuote(admin.String()), "Space Admin", name.ValueString())
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to attach space admin to space, got error: %s", err))
				return
			}
		}
	}

	if agent.IsNull() == false {
		elements := make(map[string]types.String, len(agent.Elements()))
		diags := agent.ElementsAs(ctx, &elements, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		ns := elements["default_namespace"]
		sa := elements["default_service_account"]
		agent_name := elements["name"]
		err := r.client.AddAgentToSpace(trimQuote(agent_name.String()), trimQuote(ns.String()), trimQuote(sa.String()), name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to attach agent to space, got error: %s", err))
			return
		}
	}

	tflog.Trace(ctx, "Resource Created Successful!")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueSpaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TorqueSpaceResourceModel

	// Read Terraform prior state data into the model
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

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TorqueSpaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TorqueSpaceResourceModel

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

func (r *TorqueSpaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TorqueSpaceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.DeleteSpace(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete space, got error: %s", err))
		return
	}

}

func (r *TorqueSpaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
