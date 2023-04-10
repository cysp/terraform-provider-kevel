package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	adzerk "github.com/cysp/adzerk-management-sdk-go"
)

type siteResourceModel struct {
	Id    types.Int64  `tfsdk:"id"`
	Title types.String `tfsdk:"title"`
	Url   types.String `tfsdk:"url"`
}

func (m *siteResourceModel) createRequestBody() adzerk.CreateSiteJSONRequestBody {
	return adzerk.CreateSiteJSONRequestBody{
		Title: m.Title.ValueString(),
		URL:   m.Url.ValueString(),
	}
}
func (m *siteResourceModel) updateRequestBody() adzerk.UpdateSiteJSONRequestBody {
	return adzerk.UpdateSiteJSONRequestBody{
		Id:    int32(m.Id.ValueInt64()),
		Title: m.Title.ValueString(),
		URL:   m.Url.ValueString(),
	}

}

func (m *siteResourceModel) deleteRequestBody() adzerk.UpdateSiteJSONRequestBody {
	isDeleted := true
	body := m.updateRequestBody()
	body.IsDeleted = &isDeleted
	return body
}
