package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	kevelManagementClient "github.com/cysp/adzerk-management-sdk-go"
)

func setStateWithSite(s *tfsdk.State, ctx context.Context, site *kevelManagementClient.Site) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if site == nil {
		diags.AddError("Error", "site is nil")
		return diags
	}

	SetInt64StateAttributeFromInt32Pointer(s, ctx, path.Root("id"), site.Id, &diags)
	SetStringStateAttribute(s, ctx, path.Root("title"), site.Title, &diags)
	SetStringStateAttribute(s, ctx, path.Root("url"), site.Url, &diags)

	return diags
}