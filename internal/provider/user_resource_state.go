package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	adzerk "github.com/cysp/adzerk-management-sdk-go"
)

func setStateWithUser(s *tfsdk.State, ctx context.Context, user *adzerk.User) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if user == nil {
		diags.AddError("Error", "user is nil")
		return diags
	}

	SetInt64StateAttributeFromInt32(s, ctx, path.Root("id"), user.Id, &diags)
	SetStringStateAttribute(s, ctx, path.Root("email"), user.Email, &diags)
	SetStringStateAttribute(s, ctx, path.Root("name"), user.Name, &diags)
	SetStringStateAttribute(s, ctx, path.Root("access_level"), string(user.AccessLevel), &diags)
	SetBoolStateAttribute(s, ctx, path.Root("can_access_studio"), user.CanAccessStudio, &diags)

	return diags
}
