package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	kevelManagementClient "github.com/cysp/adzerk-management-sdk-go"
)

type channelResourceModel struct {
	Id      types.Int64   `tfsdk:"id"`
	Title   types.String  `tfsdk:"title"`
	AdTypes []types.Int64 `tfsdk:"ad_types"`
}

func (m *channelResourceModel) createRequestBody() kevelManagementClient.CreateChannelJSONRequestBody {
	bodyAdTypes := make([]int32, len(m.AdTypes))
	for itemIndex, item := range m.AdTypes {
		bodyAdTypes[itemIndex] = int32(item.ValueInt64())
	}

	return kevelManagementClient.CreateChannelJSONRequestBody{
		Title:   m.Title.ValueString(),
		AdTypes: bodyAdTypes,
		Engine:  0,
	}
}
func (m *channelResourceModel) updateRequestBody() kevelManagementClient.UpdateChannelJSONRequestBody {
	bodyAdTypes := make([]int32, len(m.AdTypes))
	for itemIndex, item := range m.AdTypes {
		bodyAdTypes[itemIndex] = int32(item.ValueInt64())
	}

	return kevelManagementClient.UpdateChannelJSONRequestBody{
		Id:      int32(m.Id.ValueInt64()),
		Title:   m.Title.ValueString(),
		AdTypes: bodyAdTypes,
		Engine:  0,
	}
}
