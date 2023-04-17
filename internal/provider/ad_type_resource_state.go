package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	kevelManagementClient "github.com/cysp/adzerk-management-sdk-go"
)

func setStateWithAdType(s *tfsdk.State, ctx context.Context, adType *kevelManagementClient.AdType) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if adType == nil {
		diags.AddError("Error", "ad type is nil")
		return diags
	}

	SetInt64StateAttributeFromInt32Pointer(s, ctx, path.Root("id"), adType.Id, &diags)
	SetStringStateAttribute(s, ctx, path.Root("name"), adType.Name, &diags)
	SetInt64StateAttributeFromInt32Pointer(s, ctx, path.Root("width"), adType.Width, &diags)
	SetInt64StateAttributeFromInt32Pointer(s, ctx, path.Root("height"), adType.Height, &diags)

	return diags
}