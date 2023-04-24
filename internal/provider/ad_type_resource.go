package provider

import (
	"context"
	"sort"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

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

func (r *adTypeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ad_type"
}

func (r *adTypeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Kevel AdType",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the ad type",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
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

func (r *adTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	resp.Diagnostics.Append(setStateWithAdType(&resp.State, ctx, response.JSON200)...)
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

	adTypeIndex, adTypeFound := sort.Find(len(adTypeList.Items), func(i int) int {
		return int(int32(state.Id.ValueInt64()) - (adTypeList.Items)[i].Id)
	})

	if !adTypeFound {
		resp.Diagnostics.AddError(
			"Error Reading Kevel AdType",
			"Could not read ad type ID "+state.Id.String()+", ad type not found",
		)
		return
	}

	adType := (adTypeList.Items)[adTypeIndex]

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

	resp.Diagnostics.Append(setStateWithAdType(&resp.State, ctx, response.JSON200)...)
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
	ImportStatePassthroughInt64ID(ctx, path.Root("id"), req, resp)
}
