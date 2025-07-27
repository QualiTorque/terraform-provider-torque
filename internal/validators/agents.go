package validators

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type UseAllAgentsValidator struct{}

func (v UseAllAgentsValidator) Description(ctx context.Context) string {
	return "Ensures use_all_agents is false when agents are provided."
}

func (v UseAllAgentsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v UseAllAgentsValidator) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	useAllAgents := req.ConfigValue.ValueBool()
	var agents []types.String

	// Fetch the agents attribute
	if diags := req.Config.GetAttribute(ctx, path.Root("agents"), &agents); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Check if use_all_agents is true and agents should be empty
	if useAllAgents && len(agents) > 0 {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"Cannot specify agents when use_all_agents is true.",
		)
		return
	}

	// Check if use_all_agents is false and agents list must have at least 1 element
	if !useAllAgents && len(agents) == 0 {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"Agents list must contain at least one element when use_all_agents is false.",
		)
	}
}
