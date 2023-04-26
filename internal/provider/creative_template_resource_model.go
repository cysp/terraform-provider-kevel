package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	adzerk "github.com/cysp/adzerk-management-sdk-go"
)

type creativeTemplateResourceModel struct {
	Id          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Fields      types.List   `tfsdk:"fields"`
	Contents    types.List   `tfsdk:"contents"`
}

type creativeTemplateResourceModelField struct {
	Type        types.String `tfsdk:"type"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Variable    types.String `tfsdk:"variable"`
	Required    types.Bool   `tfsdk:"required"`
	Hidden      types.Bool   `tfsdk:"hidden"`
	AdQuery     types.Bool   `tfsdk:"ad_query"`
	Default     types.Object `tfsdk:"default"`
}

type creativeTemplateResourceModelFieldDefault struct {
	String types.String `tfsdk:"string"`
	Array  types.List   `tfsdk:"array"`
}

type creativeTemplateResourceModelContent struct {
	Type types.String `tfsdk:"type"`
	Body types.String `tfsdk:"body"`
}

// type creativeTemplateContent struct {
// 	Body *string                                           `json:"Body"`
// 	Type adzerk.CreateCreativeTemplateJSONBodyContentsType `json:"Type"`
// }

// type creativeTemplateField struct {
// 	Type        adzerk.CreateCreativeTemplateJSONBodyFieldsType `json:"Type"`
// 	Name        string                                          `json:"Name"`
// 	Description *string                                         `json:"Description"`
// 	Variable    string                                          `json:"Variable"`
// 	Required    *bool                                           `json:"Required,omitempty"`
// 	Hidden      *bool                                           `json:"Hidden,omitempty"`
// 	AdQuery     *bool                                           `json:"AdQuery"`
// 	Default     *map[string]interface{}                         `json:"Default"`
// }

func creativeTemplateResourceModelFromCreativeTemplate(ctx context.Context, creativeTemplate *adzerk.CreativeTemplate) (creativeTemplateResourceModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	modelFieldsElementDefaultAttributeTypes := map[string]attr.Type{
		"string": basetypes.StringType{},
		"array":  basetypes.ListType{ElemType: basetypes.StringType{}},
	}
	modelFieldsElementAttributeTypes := map[string]attr.Type{
		"type":        basetypes.StringType{},
		"name":        basetypes.StringType{},
		"description": basetypes.StringType{},
		"variable":    basetypes.StringType{},
		"required":    basetypes.BoolType{},
		"hidden":      basetypes.BoolType{},
		"ad_query":    basetypes.BoolType{},
		"default":     basetypes.ObjectType{AttrTypes: modelFieldsElementDefaultAttributeTypes},
	}

	modelContentsElementAttributeTypes := map[string]attr.Type{
		"type": basetypes.StringType{},
		"body": basetypes.StringType{},
	}

	modelFieldsElements := make([]creativeTemplateResourceModelField, len(creativeTemplate.Fields))
	for i, creativeTemplateField := range creativeTemplate.Fields {
		modelFieldDefault := basetypes.NewObjectNull(modelFieldsElementDefaultAttributeTypes)

		if creativeTemplateField.Default != nil {
			attributes := map[string]attr.Value{
				"string": basetypes.NewStringNull(),
				"array":  basetypes.NewListNull(basetypes.StringType{}),
			}

			if stringDefault, ok := creativeTemplateField.Default.(string); ok {
				attributes["string"] = basetypes.NewStringValue(stringDefault)
			}
			if arrayDefault, ok := creativeTemplateField.Default.([]string); ok {
				arrayDefaultValue, arrayDefaultValueDiags := basetypes.NewListValueFrom(ctx, basetypes.StringType{}, arrayDefault)
				diags.Append(arrayDefaultValueDiags...)
				attributes["array"] = arrayDefaultValue
			}

			modelFieldDefaultValue, modelFieldDefaultDiags := basetypes.NewObjectValueFrom(ctx, modelFieldsElementDefaultAttributeTypes, attributes)
			diags.Append(modelFieldDefaultDiags...)
			modelFieldDefault = modelFieldDefaultValue
		}

		modelFieldsElements[i] = creativeTemplateResourceModelField{
			Type:        basetypes.NewStringValue(creativeTemplateField.Type),
			Name:        basetypes.NewStringValue(creativeTemplateField.Name),
			Description: basetypes.NewStringValue(creativeTemplateField.Description),
			Variable:    basetypes.NewStringValue(creativeTemplateField.Variable),
			Required:    basetypes.NewBoolValue(creativeTemplateField.Required),
			Hidden:      basetypes.NewBoolValue(creativeTemplateField.Hidden),
			AdQuery:     basetypes.NewBoolValue(false),
			Default:     modelFieldDefault,
		}

	}

	modelContentsElements := make([]creativeTemplateResourceModelContent, len(creativeTemplate.Contents))

	modelFields, modelFieldsDiags := basetypes.NewListValueFrom(ctx, basetypes.ObjectType{AttrTypes: modelFieldsElementAttributeTypes}, modelFieldsElements)
	diags.Append(modelFieldsDiags...)

	modelContents, modelContentsDiags := basetypes.NewListValueFrom(ctx, basetypes.ObjectType{AttrTypes: modelContentsElementAttributeTypes}, modelContentsElements)
	diags.Append(modelContentsDiags...)

	return creativeTemplateResourceModel{
		Id:          NewInt64ValueFromInt32(creativeTemplate.Id),
		Name:        NewStringValueFromString(creativeTemplate.Name),
		Description: NewStringValueFromString(creativeTemplate.Description),
		Fields:      modelFields,
		Contents:    modelContents,
	}, diags
}

func (m *creativeTemplateResourceModel) createRequestBody(ctx context.Context) (adzerk.CreateCreativeTemplateJSONRequestBody, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	body := adzerk.CreateCreativeTemplateJSONRequestBody{
		Name:        m.Name.ValueString(),
		Description: m.Description.ValueString(),
	}

	var modelFields []creativeTemplateResourceModelField
	diags.Append(m.Fields.ElementsAs(ctx, &modelFields, false)...)

	for _, field := range modelFields {
		var fieldDefault interface{} = nil
		body.Fields = append(body.Fields, adzerk.CreateCreativeTemplateJSONBodyField{
			Type:        adzerk.CreateCreativeTemplateJSONBodyFieldsType(field.Type.ValueString()),
			Name:        field.Name.ValueString(),
			Description: field.Description.ValueStringPointer(),
			Variable:    field.Variable.ValueString(),
			Required:    field.Required.ValueBoolPointer(),
			Hidden:      field.Hidden.ValueBoolPointer(),
			AdQuery:     field.AdQuery.ValueBoolPointer(),
			Default:     &fieldDefault,
		})
	}

	return body, diags
}

func (m *creativeTemplateResourceModel) updateRequestBody() adzerk.UpdateCreativeTemplateJSONRequestBody {
	return adzerk.UpdateCreativeTemplateJSONRequestBody{
		Updates: []adzerk.CreativeTemplateUpdateOperation{
			adzerk.CreativeTemplateUpdateOperation{
				Path:  []string{"Name"},
				Value: m.Name.ValueString(),
			},
		},
	}
}

func (m *creativeTemplateResourceModel) deleteRequestBody() adzerk.UpdateCreativeTemplateJSONRequestBody {
	return adzerk.UpdateCreativeTemplateJSONRequestBody{
		Updates: []adzerk.CreativeTemplateUpdateOperation{
			adzerk.CreativeTemplateUpdateOperation{
				Path:  []string{"IsArchived"},
				Value: true,
			},
		},
	}
}
