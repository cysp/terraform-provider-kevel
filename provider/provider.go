package provider

import (
	"context"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	kevelManagementClient "github.com/cysp/terraform-provider-kevel/kevel-management-client"
)

var (
	_ provider.Provider = &kevelProvider{}
)

func New() provider.Provider {
	return &kevelProvider{}
}

type kevelProvider struct{}

type kevelProviderConfig struct {
	ApiBaseUrl types.String `tfsdk:"api_base_url"`
	ApiKey     types.String `tfsdk:"api_key"`
}

func (p *kevelProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "kevel"
}

func (p *kevelProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Kevel",
		Attributes: map[string]schema.Attribute{
			"api_base_url": schema.StringAttribute{
				Description: "The base URL of the Kevel API.",
				Optional:    true,
			},
			"api_key": schema.StringAttribute{
				Description: "Your Kevel API Key.",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *kevelProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config kevelProviderConfig

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiKeySecurityProvider, err := securityprovider.NewSecurityProviderApiKey("header", "X-Adzerk-ApiKey", config.ApiKey.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error configuring client", err.Error())
		return
	}

	var apiBaseUrl string
	if config.ApiBaseUrl.IsUnknown() {
		apiBaseUrl = "https://api.kevel.co/"
	} else {
		apiBaseUrl = config.ApiBaseUrl.ValueString()
	}

	client, err := kevelManagementClient.NewClientWithResponses(apiBaseUrl, kevelManagementClient.WithRequestEditorFn(apiKeySecurityProvider.Intercept))
	if err != nil {
		resp.Diagnostics.AddError("Error configuring client", err.Error())
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *kevelProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// NewNetworkDataSource,
	}
}

func (p *kevelProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAdTypeResource,
		NewChannelResource,
		NewChannelSiteMapResource,
		NewSiteResource,
	}
}
