package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type channelSiteMapResourceModel struct {
	ChannelId types.Int64 `tfsdk:"channel_id"`
	SiteId    types.Int64 `tfsdk:"site_id"`
	Priority  types.Int64 `tfsdk:"priority"`
}

func (m *channelSiteMapResourceModel) createRequestBody() map[string]interface{} {
	body := make(map[string]interface{}, 3)
	AddInt64ValueToMap(&body, "ChannelId", m.ChannelId)
	AddInt64ValueToMap(&body, "SiteId", m.SiteId)
	AddInt64ValueToMap(&body, "Priority", m.Priority)
	return body
}
