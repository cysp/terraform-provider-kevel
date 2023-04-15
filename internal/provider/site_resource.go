package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	kevelManagementClient "github.com/cysp/adzerk-management-sdk-go"
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
	client *kevelManagementClient.ClientWithResponses
}

type siteResourceModel struct {
	Id    types.Int64  `tfsdk:"id"`
	Title types.String `tfsdk:"title"`
	Url   types.String `tfsdk:"url"`
}

func (m *siteResourceModel) createRequestBody() map[string]interface{} {
	body := make(map[string]interface{})
	AddStringValueToMap(&body, "Title", m.Title)
	AddStringValueToMap(&body, "URL", m.Url)
	return body
}

func (m *siteResourceModel) updateRequestBody() map[string]interface{} {
	body := make(map[string]interface{})
	AddInt64ValueToMap(&body, "Id", m.Id)
	AddStringValueToMap(&body, "Title", m.Title)
	AddStringValueToMap(&body, "URL", m.Url)
	return body
}

func (m *siteResourceModel) deleteRequestBody() map[string]interface{} {
	body := m.updateRequestBody()
	body["IsDeleted"] = true
	return body
}

func setStateWithSite(s *tfsdk.State, ctx context.Context, site *kevelManagementClient.Site) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if site == nil {
		diags.AddError("Error", "site is nil")
		return diags
	}

	if site.Id != nil {
		diags.Append(s.SetAttribute(ctx, path.Root("id"), NewInt64PointerValueFromInt32(site.Id))...)
	}

	if site.Title != nil {
		diags.Append(s.SetAttribute(ctx, path.Root("title"), types.StringPointerValue(site.Title))...)
	}

	if site.Url != nil {
		diags.Append(s.SetAttribute(ctx, path.Root("url"), types.StringPointerValue(site.Url))...)

	}

	return diags
}

func (r *siteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site"
}

func (r *siteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Kevel Site",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the site",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Description: "Title of the site",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"url": schema.StringAttribute{
				Description: "URL of the site",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *siteResource) Configure(_ context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*kevelManagementClient.ClientWithResponses)
	if !ok {
		res.Diagnostics.AddError("Error", "Could not get client from provider data")
		return
	}

	r.client = client
}

func (r *siteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan siteResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.CreateSiteWithResponse(ctx, plan.createRequestBody())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating site",
			"Could not create site, unexpected error: "+err.Error(),
		)
		return
	}

	site := response.JSON200

	resp.Diagnostics.Append(setStateWithSite(&resp.State, ctx, site)...)
}

func (r *siteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state siteResourceModel

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

	site := response.JSON200

	resp.Diagnostics.Append(setStateWithSite(&resp.State, ctx, site)...)
}

func (r *siteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan siteResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.UpdateSiteWithResponse(ctx, int32(plan.Id.ValueInt64()), plan.updateRequestBody())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Kevel Site",
			"Could not update site ID "+plan.Id.String()+", unexpected error: "+err.Error(),
		)
		return
	}

	site := response.JSON200

	resp.Diagnostics.Append(setStateWithSite(&resp.State, ctx, site)...)
}

func (r *siteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state siteResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.UpdateSiteWithResponse(ctx, int32(state.Id.ValueInt64()), state.deleteRequestBody())
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
	id, err := strconv.ParseUint(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing Kevel Site",
			"Could not import site ID "+req.ID+", unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
