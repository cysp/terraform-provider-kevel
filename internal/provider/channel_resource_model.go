package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	adzerk "github.com/cysp/adzerk-management-sdk-go"
)

type channelResourceModel struct {
	Id      types.Int64  `tfsdk:"id"`
	Title   types.String `tfsdk:"title"`
	AdTypes types.List   `tfsdk:"ad_types"`
}

func (m *channelResourceModel) createRequestBody(ctx context.Context, diags *diag.Diagnostics) adzerk.CreateChannelJSONRequestBody {
	bodyAdTypes, bodyAdTypesDiags := makeRequestBodyAdTypes(ctx, m.AdTypes)
	diags.Append(bodyAdTypesDiags...)

	return adzerk.CreateChannelJSONRequestBody{
		Title:   m.Title.ValueString(),
		AdTypes: bodyAdTypes,
		Engine:  0,
	}
}

func (m *channelResourceModel) updateRequestBody(ctx context.Context, diags *diag.Diagnostics) adzerk.UpdateChannelJSONRequestBody {
	bodyAdTypes, bodyAdTypesDiags := makeRequestBodyAdTypes(ctx, m.AdTypes)
	diags.Append(bodyAdTypesDiags...)

	return adzerk.UpdateChannelJSONRequestBody{
		Id:      int32(m.Id.ValueInt64()),
		Title:   m.Title.ValueString(),
		AdTypes: bodyAdTypes,
		Engine:  0,
	}
}

func makeRequestBodyAdTypes(ctx context.Context, model types.List) ([]int32, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if model.IsNull() || model.IsUnknown() {
		return []int32{}, diags
	}

	var adTypesElements = []int64{}
	diags.Append(model.ElementsAs(ctx, &adTypesElements, false)...)
	if diags.HasError() {
		return nil, diags
	}

	bodyAdTypes := make([]int32, len(adTypesElements))
	for i, adType := range adTypesElements {
		bodyAdTypes[i] = int32(adType)
	}

	return bodyAdTypes, diags
}
