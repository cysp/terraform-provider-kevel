package provider

import "github.com/hashicorp/terraform-plugin-framework/types/basetypes"

func NewInt64PointerValueFromInt32(value *int32) basetypes.Int64Value {
	if value == nil {
		return basetypes.NewInt64Null()
	}

	return basetypes.NewInt64Value(int64(*value))
}
