package provider

import (
	"context"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

	kevelManagementClient "github.com/cysp/adzerk-management-sdk-go"
)

var (
	_ resource.Resource                = &channelSiteMapResource{}
	_ resource.ResourceWithConfigure   = &channelSiteMapResource{}
	_ resource.ResourceWithImportState = &channelSiteMapResource{}
)

func NewChannelSiteMapResource() resource.Resource {
	return &channelSiteMapResource{}
}

type channelSiteMapResource struct {
	client *kevelManagementClient.ClientWithResponses
}

func (r *channelSiteMapResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_channel_site_map"
}

func (r *channelSiteMapResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Kevel Channel Site Map",
		Attributes: map[string]schema.Attribute{
			"channel_id": schema.Int64Attribute{
				Description: "Numeric identifier of the channel",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"site_id": schema.Int64Attribute{
				Description: "Numeric identifier of the site",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"priority": schema.Int64Attribute{
				Description: "Priority of the channel site map",
				Required:    true,
			},
		},
	}
}

func (r *channelSiteMapResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*kevelManagementClient.ClientWithResponses)
	if !ok {
		resp.Diagnostics.AddError("Error", "Could not get client from provider data")
		return
	}

	r.client = client
}

func (r *channelSiteMapResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan channelSiteMapResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.CreateChannelSiteMapWithResponse(ctx, plan.createRequestBody())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating channel site map",
			"Could not create channel site map, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(setStateWithChannelSiteMap(&resp.State, ctx, response.JSON200)...)
}

func (r *channelSiteMapResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state channelSiteMapResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.GetChannelSiteMapWithResponse(ctx, int32(state.ChannelId.ValueInt64()), int32(state.SiteId.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kevel Channel Site Map",
			"Could not read channel site map "+state.ChannelId.String()+":"+state.SiteId.String()+", unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(setStateWithChannelSiteMap(&resp.State, ctx, response.JSON200)...)
}

func (r *channelSiteMapResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan channelSiteMapResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.UpdateChannelSiteMapWithResponse(ctx, plan.updateRequestBody())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Kevel Channel Site Map",
			"Could not update channel site map "+plan.ChannelId.String()+":"+plan.SiteId.String()+", unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(setStateWithChannelSiteMap(&resp.State, ctx, response.JSON200)...)
}

func (r *channelSiteMapResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state channelSiteMapResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.DeleteChannelSiteMapWithResponse(ctx, int32(state.ChannelId.ValueInt64()), int32(state.SiteId.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Kevel Channel Site Map",
			"Could not delete channel site map "+state.ChannelId.String()+":"+state.SiteId.String()+", unexpected error: "+err.Error(),
		)
		return
	}

	if response.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Error Deleting Kevel Site",
			"Could not delete channel site map "+state.ChannelId.String()+":"+state.SiteId.String()+", unexpected status code: "+strconv.Itoa(response.StatusCode()),
		)
		return
	}
}

var importChannelSiteMapResourceIdRegExp = regexp.MustCompile("^([0-9]+):([0-9]+)$")

func (r *channelSiteMapResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	matches := importChannelSiteMapResourceIdRegExp.FindStringSubmatch(req.ID)
	if len(matches) != 3 {
		resp.Diagnostics.AddError(
			"Error importing Kevel Channel Site Map",
			"Could not import channel site map, error parsing identifier: "+req.ID,
		)
		return
	}

	channelId, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing Kevel Channel Site Map",
			"Could not import channel site map ID "+req.ID+", unexpected error: "+err.Error(),
		)
		return
	}

	siteId, err := strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing Kevel Channel Site Map",
			"Could not import channel site map ID "+req.ID+", unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("channel_id"), channelId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("site_id"), siteId)...)
}
