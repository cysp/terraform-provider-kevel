package provider

import "github.com/hashicorp/terraform-plugin-framework/types/basetypes"

func NewInt64PointerValueFromInt32(value *int32) basetypes.Int64Value {
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
