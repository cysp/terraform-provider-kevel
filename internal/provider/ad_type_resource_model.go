package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	adzerk "github.com/cysp/adzerk-management-sdk-go"
)

type adTypeResourceModel struct {
	Id     types.Int64  `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Width  types.Int64  `tfsdk:"width"`
	Height types.Int64  `tfsdk:"height"`
}

func (m *adTypeResourceModel) createRequestBody() adzerk.CreateAdTypeJSONRequestBody {
	return adzerk.CreateAdTypeJSONRequestBody{
		Width:  int32(m.Width.ValueInt64()),
		Height: int32(m.Height.ValueInt64()),
		Name:   m.Name.ValueStringPointer(),
	}
}
