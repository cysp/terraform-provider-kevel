package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	kevelManagementClient "github.com/cysp/adzerk-management-sdk-go"
)

type siteResourceModel struct {
	Id    types.Int64  `tfsdk:"id"`
	Title types.String `tfsdk:"title"`
	Url   types.String `tfsdk:"url"`
}

func (m *siteResourceModel) createRequestBody() kevelManagementClient.CreateSiteJSONRequestBody {
	return kevelManagementClient.CreateSiteJSONRequestBody{
		Title: m.Title.ValueString(),
		URL:   m.Url.ValueString(),
	}
}
func (m *siteResourceModel) updateRequestBody() kevelManagementClient.UpdateSiteJSONRequestBody {
	return kevelManagementClient.UpdateSiteJSONRequestBody{
		Id:    int32(m.Id.ValueInt64()),
		Title: m.Title.ValueString(),
		URL:   m.Url.ValueString(),
	}

}

func (m *siteResourceModel) deleteRequestBody() kevelManagementClient.UpdateSiteJSONRequestBody {
	isDeleted := true
	body := m.updateRequestBody()
	body.IsDeleted = &isDeleted
	return body
}
