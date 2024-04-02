package resource_site

import (
	adzerk "github.com/cysp/adzerk-management-sdk-go"
)

func (m *SiteModel) CreateRequestBody() adzerk.CreateSiteJSONRequestBody {
	return adzerk.CreateSiteJSONRequestBody{
		Title: m.Title.ValueString(),
		URL:   m.Url.ValueString(),
	}
}
func (m *SiteModel) UpdateRequestBody() adzerk.UpdateSiteJSONRequestBody {
	return adzerk.UpdateSiteJSONRequestBody{
		Id:    int32(m.Id.ValueInt64()),
		Title: m.Title.ValueString(),
		URL:   m.Url.ValueString(),
	}

}

func (m *SiteModel) DeleteRequestBody() adzerk.UpdateSiteJSONRequestBody {
	isDeleted := true
	body := m.UpdateRequestBody()
	body.IsDeleted = &isDeleted
	return body
}
