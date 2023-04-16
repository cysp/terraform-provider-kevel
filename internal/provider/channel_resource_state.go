package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	kevelManagementClient "github.com/cysp/adzerk-management-sdk-go"
)

func setStateWithChannel(s *tfsdk.State, ctx context.Context, channel *kevelManagementClient.Channel) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if channel == nil {
		diags.AddError("Error", "channel is nil")
		return diags
	}

	if channel.Id != nil {
		diags.Append(s.SetAttribute(ctx, path.Root("id"), NewInt64ValueFromInt32Pointer(channel.Id))...)
	}

	if channel.Title != nil {
		diags.Append(s.SetAttribute(ctx, path.Root("title"), types.StringPointerValue(channel.Title))...)
	}

	if channel.AdTypes != nil {
		stateAdTypes := make([]basetypes.Int64Value, len(*channel.AdTypes))
		for itemIndex, item := range *channel.AdTypes {
			stateAdTypes[itemIndex] = types.Int64Value(int64(item))
		}

		diags.Append(s.SetAttribute(ctx, path.Root("ad_types"), stateAdTypes)...)
	}

	return diags
}
