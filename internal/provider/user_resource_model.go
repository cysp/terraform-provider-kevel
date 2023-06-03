package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	adzerk "github.com/cysp/adzerk-management-sdk-go"
)

type userResourceModel struct {
	Id              types.Int64  `tfsdk:"id"`
	Email           types.String `tfsdk:"email"`
	Name            types.String `tfsdk:"name"`
	AccessLevel     types.String `tfsdk:"access_level"`
	CanAccessStudio types.Bool   `tfsdk:"can_access_studio"`
}

func (m *userResourceModel) createRequestBody() adzerk.CreateUserJSONRequestBody {
	return adzerk.CreateUserJSONRequestBody{
		Email:           m.Email.ValueString(),
		Name:            m.Name.ValueStringPointer(),
		AccessLevel:     (*adzerk.AccessLevel)(m.AccessLevel.ValueStringPointer()),
		CanAccessStudio: m.CanAccessStudio.ValueBoolPointer(),
	}
}
func (m *userResourceModel) updateRequestBody() adzerk.UpdateUserJSONRequestBody {
	return adzerk.UpdateUserJSONRequestBody{
		Id:              int32(m.Id.ValueInt64()),
		Email:           m.Email.ValueString(),
		Name:            m.Name.ValueString(),
		AccessLevel:     (*adzerk.AccessLevel)(m.AccessLevel.ValueStringPointer()),
		CanAccessStudio: m.CanAccessStudio.ValueBoolPointer(),
	}

}

// func (m *userResourceModel) deleteRequestBody() adzerk.UpdateUserJSONRequestBody {
// 	isDeleted := true
// 	body := m.updateRequestBody()
// 	body.IsDeleted = &isDeleted
// 	return body
// }
