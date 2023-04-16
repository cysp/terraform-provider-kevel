package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type siteResourceModel struct {
	Id    types.Int64  `tfsdk:"id"`
	Title types.String `tfsdk:"title"`
	Url   types.String `tfsdk:"url"`
}

func (m *siteResourceModel) createOrUpdateRequestBody() map[string]interface{} {
	body := make(map[string]interface{})
	AddInt64ValueToMap(&body, "Id", m.Id)
	AddStringValueToMap(&body, "Title", m.Title)
	AddStringValueToMap(&body, "URL", m.Url)
	return body
}

func (m *siteResourceModel) deleteRequestBody() map[string]interface{} {
	body := m.createOrUpdateRequestBody()
	body["IsDeleted"] = true
	return body
}
