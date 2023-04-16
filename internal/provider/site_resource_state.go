package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	kevelManagementClient "github.com/cysp/adzerk-management-sdk-go"
)

func setStateWithSite(s *tfsdk.State, ctx context.Context, site *kevelManagementClient.Site) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if site == nil {
		diags.AddError("Error", "site is nil")
		return diags
	}

	if site.Id != nil {
		diags.Append(s.SetAttribute(ctx, path.Root("id"), NewInt64ValueFromInt32Pointer(site.Id))...)
	}

	if site.Title != nil {
		diags.Append(s.SetAttribute(ctx, path.Root("title"), types.StringPointerValue(site.Title))...)
	}

	if site.Url != nil {
		diags.Append(s.SetAttribute(ctx, path.Root("url"), types.StringPointerValue(site.Url))...)

	}

	return diags
}
