package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	adzerk "github.com/cysp/adzerk-management-sdk-go"
	"github.com/cysp/terraform-provider-kevel/internal/provider/resource_site"
)

var (
	_ resource.Resource                = &siteResource{}
	_ resource.ResourceWithConfigure   = &siteResource{}
	_ resource.ResourceWithImportState = &siteResource{}
)

func NewSiteResource() resource.Resource {
	return &siteResource{}
}

type siteResource struct {
	client *adzerk.ClientWithResponses
}

func (r *siteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site"
}

func (r *siteResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_site.SiteResourceSchema(ctx)
	resp.Schema.Description = "Kevel Site"
}

func (r *siteResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*adzerk.ClientWithResponses)
	if !ok {
		resp.Diagnostics.AddError("Error", "Could not get client from provider data")
		return
	}

	r.client = client
}

func (r *siteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resource_site.SiteModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.CreateSiteWithResponse(ctx, plan.CreateRequestBody())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating site",
			"Could not create site, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(setStateWithSite(&resp.State, ctx, response.JSON200)...)
}

func (r *siteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resource_site.SiteModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.GetSiteWithResponse(ctx, int32(state.Id.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kevel Site",
			"Could not read site ID "+state.Id.String()+", unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(setStateWithSite(&resp.State, ctx, response.JSON200)...)
}

func (r *siteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan resource_site.SiteModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.UpdateSiteWithResponse(ctx, int32(plan.Id.ValueInt64()), plan.UpdateRequestBody())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Kevel Site",
			"Could not update site ID "+plan.Id.String()+", unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(setStateWithSite(&resp.State, ctx, response.JSON200)...)
}

func (r *siteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resource_site.SiteModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.UpdateSiteWithResponse(ctx, int32(state.Id.ValueInt64()), state.DeleteRequestBody())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Kevel Site",
			"Could not delete site ID "+state.Id.String()+", unexpected error: "+err.Error(),
		)
		return
	}

	if response.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Error Deleting Kevel Site",
			"Could not delete site ID "+state.Id.String()+", unexpected status code: "+strconv.Itoa(response.StatusCode()),
		)
		return
	}
}

func (r *siteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ImportStatePassthroughInt64ID(ctx, path.Root("id"), req, resp)
}
