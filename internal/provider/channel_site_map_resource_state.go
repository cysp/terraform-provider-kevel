package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	kevelManagementClient "github.com/cysp/adzerk-management-sdk-go"
)

func setStateWithChannelSiteMap(s *tfsdk.State, ctx context.Context, channelSiteMap *kevelManagementClient.ChannelSiteMap) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if channelSiteMap == nil {
		diags.AddError("Error", "channel site map is nil")
		return diags
	}

	if channelSiteMap.ChannelId != nil {
		diags.Append(s.SetAttribute(ctx, path.Root("channel_id"), NewInt64ValueFromInt32Pointer(channelSiteMap.ChannelId))...)
	}
	if channelSiteMap.SiteId != nil {
		diags.Append(s.SetAttribute(ctx, path.Root("site_id"), NewInt64ValueFromInt32Pointer(channelSiteMap.SiteId))...)
	}
	if channelSiteMap.Priority != nil {
		diags.Append(s.SetAttribute(ctx, path.Root("priority"), NewInt64ValueFromInt32Pointer(channelSiteMap.Priority))...)
	}

	return diags
}
