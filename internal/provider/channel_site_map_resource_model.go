package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	adzerk "github.com/cysp/adzerk-management-sdk-go"
)

type channelSiteMapResourceModel struct {
	Id        types.String `tfsdk:"id"`
	ChannelId types.Int64  `tfsdk:"channel_id"`
	SiteId    types.Int64  `tfsdk:"site_id"`
	Priority  types.Int64  `tfsdk:"priority"`
}

func (m *channelSiteMapResourceModel) createRequestBody() adzerk.CreateChannelSiteMapJSONRequestBody {
	return adzerk.CreateChannelSiteMapJSONRequestBody{
		ChannelId: int32(m.ChannelId.ValueInt64()),
		SiteId:    int32(m.SiteId.ValueInt64()),
		Priority:  int32(m.Priority.ValueInt64()),
	}
}

func (m *channelSiteMapResourceModel) updateRequestBody() adzerk.UpdateChannelSiteMapJSONRequestBody {
	return adzerk.UpdateChannelSiteMapJSONRequestBody{
		ChannelId: int32(m.ChannelId.ValueInt64()),
		SiteId:    int32(m.SiteId.ValueInt64()),
		Priority:  int32(m.Priority.ValueInt64()),
	}
}
