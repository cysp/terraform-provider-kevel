package provider

import (
	"context"
	"sort"
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
	_ resource.Resource                = &adTypeResource{}
	_ resource.ResourceWithConfigure   = &adTypeResource{}
	_ resource.ResourceWithImportState = &adTypeResource{}
)

func NewAdTypeResource() resource.Resource {
	return &adTypeResource{}
}

type adTypeResource struct {
	client *kevelManagementClient.ClientWithResponses
}

type adTypeResourceModel struct {
	Id     types.Int64  `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Width  types.Int64  `tfsdk:"width"`
	Height types.Int64  `tfsdk:"height"`
}

func (adType *adTypeResourceModel) createRequestBody() map[string]interface{} {
	body := make(map[string]interface{})
	AddInt64ValueToMap(&body, "Id", adType.Id)
	AddStringValueToMap(&body, "Name", adType.Name)
	AddInt64ValueToMap(&body, "Width", adType.Width)
	AddInt64ValueToMap(&body, "Height", adType.Height)
	return body
}

func setStateWithAdType(s *tfsdk.State, ctx context.Context, adType *kevelManagementClient.AdType) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if adType == nil {
		diags.AddError("Error", "ad type is nil")
		return diags
	}

	if adType.Id != nil {
		diags.Append(s.SetAttribute(ctx, path.Root("id"), NewInt64PointerValueFromInt32(adType.Id))...)
	}

	if adType.Name != nil {
		diags.Append(s.SetAttribute(ctx, path.Root("name"), types.StringPointerValue(adType.Name))...)
	}

	if adType.Width != nil {
		diags.Append(s.SetAttribute(ctx, path.Root("width"), NewInt64PointerValueFromInt32(adType.Width))...)
	}

	if adType.Height != nil {
		diags.Append(s.SetAttribute(ctx, path.Root("height"), NewInt64PointerValueFromInt32(adType.Height))...)
	}

	return diags
}

func (r *adTypeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ad_type"
}

func (r *adTypeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Kevel AdType",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:   "Numeric identifier of the ad type",
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{},
			},
			"name": schema.StringAttribute{
				Description: "Name of the ad type",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"width": schema.Int64Attribute{
				Description: "Width of the ad type",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"height": schema.Int64Attribute{
				Description: "Height of the ad type",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *adTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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

func (r *adTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan adTypeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.CreateAdTypeWithResponse(ctx, plan.createRequestBody())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ad type",
			"Could not create ad type, unexpected error: "+err.Error(),
		)
		return
	}

	adType := response.JSON200

	resp.Diagnostics.Append(setStateWithAdType(&resp.State, ctx, adType)...)
}

func (r *adTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state adTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.ListAdTypesWithResponse(ctx, &kevelManagementClient.ListAdTypesParams{})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kevel AdType",
			"Could not read ad type ID "+state.Id.String()+", unexpected error: "+err.Error(),
		)
		return
	}

	adTypeList := response.JSON200

	adTypeIndex, adTypeFound := sort.Find(len(*adTypeList.Items), func(i int) int {
		return int(int32(state.Id.ValueInt64()) - *(*adTypeList.Items)[i].Id)
	})

	if !adTypeFound {
		resp.Diagnostics.AddError(
			"Error Reading Kevel AdType",
			"Could not read ad type ID "+state.Id.String()+", ad type not found",
		)
		return
	}

	adType := (*adTypeList.Items)[adTypeIndex]

	resp.Diagnostics.Append(setStateWithAdType(&resp.State, ctx, &adType)...)
}

func (r *adTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan adTypeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.CreateAdTypeWithResponse(ctx, plan.createRequestBody())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Kevel AdType",
			"Could not update ad type ID "+plan.Id.String()+", unexpected error: "+err.Error(),
		)
		return
	}

	adType := response.JSON200

	resp.Diagnostics.Append(setStateWithAdType(&resp.State, ctx, adType)...)
}

func (r *adTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state adTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.DeleteAdTypeWithResponse(ctx, int32(state.Id.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Kevel AdType",
			"Could not delete ad type ID "+state.Id.String()+", unexpected error: "+err.Error(),
		)
		return
	}

	if response.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Error Deleting Kevel AdType",
			"Could not delete ad type ID "+state.Id.String()+", unexpected status code: "+strconv.Itoa(response.StatusCode()),
		)
		return
	}
}

func (r *adTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseUint(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing Kevel AdType",
			"Could not import ad type ID "+req.ID+", unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
