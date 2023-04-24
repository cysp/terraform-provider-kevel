package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func NewInt64ValueFromInt32Pointer(value *int32) basetypes.Int64Value {
	if value == nil {
		return basetypes.NewInt64Null()
	}

	return basetypes.NewInt64Value(int64(*value))
}

func AddInt64ValueToMap(m *map[string]interface{}, key string, value basetypes.Int64Value) {
	if value.IsUnknown() {
		return
	}

	if value.IsNull() {
		(*m)[key] = nil
	} else {
		(*m)[key] = value.ValueInt64()
	}
}

func AddStringValueToMap(m *map[string]interface{}, key string, value basetypes.StringValue) {
	if value.IsUnknown() {
		return
	}

	if value.IsNull() {
		(*m)[key] = nil
	} else {
		(*m)[key] = value.ValueString()
	}
}

func SetInt64StateAttributeFromInt32(s *tfsdk.State, ctx context.Context, path path.Path, value int32, diags *diag.Diagnostics) {
	diags.Append(s.SetAttribute(ctx, path, int64(value))...)
}

func SetInt64StateAttributeFromInt32Pointer(s *tfsdk.State, ctx context.Context, path path.Path, value *int32, diags *diag.Diagnostics) {
	diags.Append(s.SetAttribute(ctx, path, NewInt64ValueFromInt32Pointer(value))...)
}

func SetStringStateAttribute(s *tfsdk.State, ctx context.Context, path path.Path, value string, diags *diag.Diagnostics) {
	diags.Append(s.SetAttribute(ctx, path, (value))...)
}
func SetStringStateAttributeFromPointer(s *tfsdk.State, ctx context.Context, path path.Path, value *string, diags *diag.Diagnostics) {
	diags.Append(s.SetAttribute(ctx, path, types.StringPointerValue(value))...)
}

func ImportStatePassthroughInt64ID(ctx context.Context, attrPath path.Path, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Resource Import Passthrough Invalid ID",
			"Could not parse ID \""+req.ID+"\" as integer: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, attrPath, id)...)
}
