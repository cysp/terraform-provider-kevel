package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	adzerk "github.com/cysp/adzerk-management-sdk-go"
)

func setStateWithCreativeTemplate(s *tfsdk.State, ctx context.Context, creativeTemplate *adzerk.CreativeTemplate) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if creativeTemplate == nil {
		diags.AddError("Error", "creativeTemplate is nil")
		return diags
	}

	model, modelDiags := creativeTemplateResourceModelFromCreativeTemplate(ctx, creativeTemplate)
	diags.Append(modelDiags...)

	s.Set(ctx, model)

	// SetInt64StateAttributeFromInt32(s, ctx, path.Root("id"), creativeTemplate.Id, &diags)
	// SetStringStateAttribute(s, ctx, path.Root("name"), creativeTemplate.Name, &diags)
	// SetStringStateAttribute(s, ctx, path.Root("description"), creativeTemplate.Description, &diags)

	// stateFields := make([]creativeTemplateResourceModelField, len(creativeTemplate.Fields))
	// for i, creativeTemplateField := range creativeTemplate.Fields {
	// 	var stateFieldDefault basetypes.ObjectValue

	// 	if creativeTemplateField.Default != nil {
	// 		attributeTypes := map[string]attr.Type{
	// 			"string": types.TypeString(),
	// 			"array":  types.TypeList,
	// 		}
	// 		attributes := map[string]attr.Value{
	// 			"string": basetypes.NewStringNull(),
	// 			"array":  basetypes.NewListNull(),
	// 		}

	// 		if stringDefault, ok := creativeTemplateField.Default.(string); ok {
	// 			attributes["string"] = basetypes.NewStringValue(stringDefault)
	// 		}
	// 		if arrayDefault, ok := creativeTemplateField.Default.([]string); ok {
	// 			attributes["array"] = basetypes.NewListValue(types.String, (arrayDefault))
	// 		}
	// 		// := creativeTemplateResourceModelFieldDefault{}
	// 	}

	// 	stateFields[i] = creativeTemplateResourceModelField{
	// 		Type:        basetypes.NewStringValue(creativeTemplateField.Type),
	// 		Name:        basetypes.NewStringValue(creativeTemplateField.Name),
	// 		Description: basetypes.NewStringValue(creativeTemplateField.Description),
	// 		Variable:    basetypes.NewStringValue(creativeTemplateField.Variable),
	// 		Required:    basetypes.NewBoolValue(creativeTemplateField.Required),
	// 		Hidden:      basetypes.NewBoolValue(creativeTemplateField.Hidden),
	// 		// AdQuery: 	basetypes.NewBoolValue(creativeTemplateField.),
	// 		Default: stateFieldDefault,
	// 	}

	// }
	// diags.Append(s.SetAttribute(ctx, path.Root("fields"), (creativeTemplate.Fields))...)
	// diags.Append(s.SetAttribute(ctx, path.Root("contents"), (creativeTemplate.Contents))...)

	return diags
}
