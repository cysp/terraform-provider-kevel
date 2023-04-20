package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	kevelManagementClient "github.com/cysp/adzerk-management-sdk-go"
)

type channelSiteMapResourceModel struct {
	ChannelId types.Int64 `tfsdk:"channel_id"`
	SiteId    types.Int64 `tfsdk:"site_id"`
	Priority  types.Int64 `tfsdk:"priority"`
}

func (m *channelSiteMapResourceModel) createRequestBody() kevelManagementClient.CreateChannelSiteMapJSONRequestBody {
	priority := int32(m.Priority.ValueInt64())

	return kevelManagementClient.CreateChannelSiteMapJSONRequestBody{
		ChannelId: int32(m.ChannelId.ValueInt64()),
		SiteId:    int32(m.SiteId.ValueInt64()),
		Priority:  &priority,
	}
}

func (m *channelSiteMapResourceModel) updateRequestBody() kevelManagementClient.UpdateChannelSiteMapJSONRequestBody {
	return kevelManagementClient.UpdateChannelSiteMapJSONRequestBody{
		ChannelId: int32(m.ChannelId.ValueInt64()),
		SiteId:    int32(m.SiteId.ValueInt64()),
		Priority:  int32(m.Priority.ValueInt64()),
	}
}
