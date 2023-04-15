package provider

import (
	"context"
	"os"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	adzerkManagementSdk "github.com/cysp/adzerk-management-sdk-go"
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

// KevelProviderModel describes the provider data model.
type KevelProviderModel struct {
	ApiBaseUrl types.String `tfsdk:"api_base_url"`
	ApiKey     types.String `tfsdk:"api_key"`
}

func (p *KevelProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "kevel"
	resp.Version = p.version
}

func (p *KevelProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_base_url": schema.StringAttribute{
				Description: "The base URL of the Kevel API.",
				Optional:    true,
			},
			"api_key": schema.StringAttribute{
				Description: "Your Kevel API Key.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *KevelProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data KevelProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var apiBaseUrl string
	if !data.ApiBaseUrl.IsNull() {
		apiBaseUrl = data.ApiBaseUrl.ValueString()
	} else {
		apiBaseUrl = "https://api.kevel.co/"
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

	client, err := adzerkManagementSdk.NewClientWithResponses(apiBaseUrl, adzerkManagementSdk.WithRequestEditorFn(apiKeySecurityProvider.Intercept))
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
