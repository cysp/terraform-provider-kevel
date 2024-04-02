// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package provider_kevel

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
)

func KevelProviderSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_base_url": schema.StringAttribute{
				Optional:            true,
				Description:         "The base URL of the Kevel API. This can also be set via the KEVEL_API_BASE_URL environment variable.",
				MarkdownDescription: "The base URL of the Kevel API. This can also be set via the KEVEL_API_BASE_URL environment variable.",
			},
			"api_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				Description:         "Your Kevel API Key. This can also be set via the KEVEL_API_KEY environment variable.",
				MarkdownDescription: "Your Kevel API Key. This can also be set via the KEVEL_API_KEY environment variable.",
			},
		},
	}
}

type KevelModel struct {
	ApiBaseUrl types.String `tfsdk:"api_base_url"`
	ApiKey     types.String `tfsdk:"api_key"`
}
