package provider

import (
	"context"
	"os"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	adzerk "github.com/cysp/adzerk-management-sdk-go"
	"github.com/cysp/terraform-provider-kevel/internal/provider/provider_kevel"
)

// Ensure KevelProvider satisfies various provider interfaces.
var (
	_ provider.Provider = &KevelProvider{}
)

// KevelProvider defines the provider implementation.
type KevelProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

func (p *KevelProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "kevel"
	resp.Version = p.version
}

func (p *KevelProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = provider_kevel.KevelProviderSchema(ctx)
	resp.Schema.Description = "The \"kevel\" provider allows the configuration of inventory items within the Kevel ad server platform."
	resp.Schema.MarkdownDescription = "The \"kevel\" provider allows the configuration of inventory items within the [Kevel](https://www.kevel.com) ad server platform."
}

func (p *KevelProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data provider_kevel.KevelModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiBaseUrl string
	if !data.ApiBaseUrl.IsNull() {
		apiBaseUrl = data.ApiBaseUrl.ValueString()
	} else {
		kevelApiBaseUrl, found := os.LookupEnv("KEVEL_API_BASE_URL")
		if found {
			apiBaseUrl = kevelApiBaseUrl
		} else {
			apiBaseUrl = "https://api.kevel.co/"
		}
	}

	if apiBaseUrl == "" {
		resp.Diagnostics.AddError("Error configuring client", "No API base URL provided")
		return
	}

	var apiKey string
	if !data.ApiKey.IsNull() {
		apiKey = data.ApiKey.ValueString()
	} else {
		apiKey = os.Getenv("KEVEL_API_KEY")
	}

	if apiKey == "" {
		resp.Diagnostics.AddError("Error configuring client", "No API key provided")
		return
	}

	apiKeySecurityProvider, err := securityprovider.NewSecurityProviderApiKey("header", "X-Adzerk-ApiKey", apiKey)
	if err != nil {
		resp.Diagnostics.AddError("Error configuring client", err.Error())
		return
	}

	client, err := adzerk.NewClientWithResponses(apiBaseUrl, adzerk.WithRequestEditorFn(apiKeySecurityProvider.Intercept))
	if err != nil {
		resp.Diagnostics.AddError("Error configuring client", err.Error())
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *KevelProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAdTypeResource,
		NewChannelResource,
		NewChannelSiteMapResource,
		NewSiteResource,
	}
}

func (p *KevelProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &KevelProvider{
			version: version,
		}
	}
}
