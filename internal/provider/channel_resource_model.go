package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type channelResourceModel struct {
	Id      types.Int64   `tfsdk:"id"`
	Title   types.String  `tfsdk:"title"`
	AdTypes []types.Int64 `tfsdk:"ad_types"`
}

func (m *channelResourceModel) createOrUpdateRequestBody() map[string]interface{} {
	bodyAdTypes := make([]int32, len(m.AdTypes))
	for itemIndex, item := range m.AdTypes {
		bodyAdTypes[itemIndex] = int32(item.ValueInt64())
	}

	body := make(map[string]interface{})
	AddInt64ValueToMap(&body, "Id", m.Id)
	AddStringValueToMap(&body, "Title", m.Title)
	body["AdTypes"] = bodyAdTypes
	body["Engine"] = 0
	return body
}
