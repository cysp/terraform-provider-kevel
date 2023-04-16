package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type adTypeResourceModel struct {
	Id     types.Int64  `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Width  types.Int64  `tfsdk:"width"`
	Height types.Int64  `tfsdk:"height"`
}

func (m *adTypeResourceModel) createRequestBody() map[string]interface{} {
	body := make(map[string]interface{})
	AddStringValueToMap(&body, "Name", m.Name)
	AddInt64ValueToMap(&body, "Width", m.Width)
	AddInt64ValueToMap(&body, "Height", m.Height)
	return body
}
