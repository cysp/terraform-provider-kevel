package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	kevelManagementClient "github.com/cysp/terraform-provider-kevel/kevel-management-client"
)

var (
	_ resource.Resource                = &channelResource{}
	_ resource.ResourceWithConfigure   = &channelResource{}
	_ resource.ResourceWithImportState = &channelResource{}
)

func NewChannelResource() resource.Resource {
	return &channelResource{}
}

type channelResource struct {
	client *kevelManagementClient.ClientWithResponses
}

type channelResourceModel struct {
	Id      types.Int64   `tfsdk:"id"`
	Title   types.String  `tfsdk:"title"`
	AdTypes []types.Int64 `tfsdk:"ad_types"`
}

func (channel *channelResourceModel) createRequestBody() kevelManagementClient.CreateChannelJSONRequestBody {
	bodyAdTypes := make([]int32, len(channel.AdTypes))
	for itemIndex, item := range channel.AdTypes {
		bodyAdTypes[itemIndex] = int32(item.ValueInt64())
	}
	return kevelManagementClient.CreateChannelJSONRequestBody{
		Title:   channel.Title.ValueString(),
		AdTypes: bodyAdTypes,
		Engine:  0,
	}
}
func (channel *channelResourceModel) updateRequestBody() kevelManagementClient.UpdateChannelJSONRequestBody {
	bodyAdTypes := make([]int32, len(channel.AdTypes))
	for itemIndex, item := range channel.AdTypes {
		bodyAdTypes[itemIndex] = int32(item.ValueInt64())
	}
	return kevelManagementClient.UpdateChannelJSONRequestBody{
		Id:      int32(channel.Id.ValueInt64()),
		Title:   channel.Title.ValueString(),
		AdTypes: bodyAdTypes,
		Engine:  0,
	}
}

func setStateWithChannel(s *tfsdk.State, ctx context.Context, channel *kevelManagementClient.Channel) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if channel == nil {
		diags.AddError("Error", "channel is nil")
		return diags
	}

	if channel.Id != nil {
		diags.Append(s.SetAttribute(ctx, path.Root("id"), NewInt64PointerValueFromInt32(channel.Id))...)
	}

	if channel.Title != nil {
		diags.Append(s.SetAttribute(ctx, path.Root("title"), types.StringPointerValue(channel.Title))...)
	}

	if channel.AdTypes != nil {
		stateAdTypes := make([]basetypes.Int64Value, len(*channel.AdTypes))
		for itemIndex, item := range *channel.AdTypes {
			stateAdTypes[itemIndex] = types.Int64Value(int64(item))
		}

		diags.Append(s.SetAttribute(ctx, path.Root("ad_types"), stateAdTypes)...)
	}

	return diags
}

func (r *channelResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_channel"
}

func (r *channelResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Kevel Channel",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the channel",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Description: "Title of the channel",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ad_types": schema.ListAttribute{
				Description: "List of ad types",
				ElementType: types.Int64Type,
				Optional:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *channelResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*kevelManagementClient.ClientWithResponses)
}

func (r *channelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan channelResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.CreateChannelWithResponse(ctx, plan.createRequestBody())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating channel",
			"Could not create channel, unexpected error: "+err.Error(),
		)
		return
	}

	channel := response.JSON200

	resp.Diagnostics.Append(setStateWithChannel(&resp.State, ctx, channel)...)
}

func (r *channelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state channelResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.GetChannelWithResponse(ctx, int32(state.Id.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kevel Channel",
			"Could not read channel ID "+state.Id.String()+", unexpected error: "+err.Error(),
		)
		return
	}

	channel := response.JSON200

	resp.Diagnostics.Append(setStateWithChannel(&resp.State, ctx, channel)...)
}

func (r *channelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan channelResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.UpdateChannelWithResponse(ctx, int32(plan.Id.ValueInt64()), plan.updateRequestBody())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Kevel Channel",
			"Could not update channel ID "+plan.Id.String()+", unexpected error: "+err.Error(),
		)
		return
	}

	channel := response.JSON200

	resp.Diagnostics.Append(setStateWithChannel(&resp.State, ctx, channel)...)
}

func (r *channelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state channelResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.DeleteChannelWithResponse(ctx, int32(state.Id.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Kevel Channel",
			"Could not delete channel ID "+state.Id.String()+", unexpected error: "+err.Error(),
		)
		return
	}

	if response.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Error Deleting Kevel Channel",
			"Could not delete channel ID "+state.Id.String()+", unexpected status code: "+strconv.Itoa(response.StatusCode()),
		)
		return
	}
}

func (r *channelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseUint(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing Kevel Channel",
			"Could not import channel ID "+req.ID+", unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
