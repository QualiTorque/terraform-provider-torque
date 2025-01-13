package common

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type KeyValuePairModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}
