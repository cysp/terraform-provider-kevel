package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	adzerk "github.com/cysp/adzerk-management-sdk-go"
)

func setStateWithChannelSiteMap(s *tfsdk.State, ctx context.Context, channelSiteMap *adzerk.ChannelSiteMap) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if channelSiteMap == nil {
		diags.AddError("Error", "channel site map is nil")
		return diags
	}

	id := fmt.Sprintf("%d:%d", channelSiteMap.ChannelId, channelSiteMap.SiteId)
	SetStringStateAttribute(s, ctx, path.Root("id"), id, &diags)

	SetInt64StateAttributeFromInt32(s, ctx, path.Root("channel_id"), channelSiteMap.ChannelId, &diags)
	SetInt64StateAttributeFromInt32(s, ctx, path.Root("site_id"), channelSiteMap.SiteId, &diags)
	SetInt64StateAttributeFromInt32Pointer(s, ctx, path.Root("priority"), channelSiteMap.Priority, &diags)

	return diags
}
