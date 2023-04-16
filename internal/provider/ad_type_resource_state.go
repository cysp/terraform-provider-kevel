package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	kevelManagementClient "github.com/cysp/adzerk-management-sdk-go"
)

func setStateWithAdType(s *tfsdk.State, ctx context.Context, adType *kevelManagementClient.AdType) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if adType == nil {
		diags.AddError("Error", "ad type is nil")
		return diags
	}

	if adType.Id != nil {
		diags.Append(s.SetAttribute(ctx, path.Root("id"), NewInt64ValueFromInt32Pointer(adType.Id))...)
	}

	if adType.Name != nil {
		diags.Append(s.SetAttribute(ctx, path.Root("name"), types.StringPointerValue(adType.Name))...)
	}

	if adType.Width != nil {
		diags.Append(s.SetAttribute(ctx, path.Root("width"), NewInt64ValueFromInt32Pointer(adType.Width))...)
	}

	if adType.Height != nil {
		diags.Append(s.SetAttribute(ctx, path.Root("height"), NewInt64ValueFromInt32Pointer(adType.Height))...)
	}

	return diags
}
