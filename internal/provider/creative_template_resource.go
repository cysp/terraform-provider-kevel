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
	_ resource.Resource                = &CreativeTemplateResource{}
	_ resource.ResourceWithConfigure   = &CreativeTemplateResource{}
	_ resource.ResourceWithImportState = &CreativeTemplateResource{}
)

func NewCreativeTemplateResource() resource.Resource {
	return &CreativeTemplateResource{}
}

type CreativeTemplateResource struct {
	client *kevelManagementClient.ClientWithResponses
}

type CreativeTemplateResourceModel struct {
	Id          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Contents    types.List   `json:"Contents,omitempty"`
	Fields      types.List   `json:"Fields,omitempty"`
}

type CreativeTemplateResourceModelContents struct {
	Type types.String `tfsdk:"type"`
	Body types.String `tfsdk:"body"`
}

type CreativeTemplateResourceModelField struct {
	Type        types.String `tfsdk:"type"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Hidden      types.Bool   `tfsdk:"hidden"`
	Required    types.Bool   `tfsdk:"required"`
	Variable    types.String `tfsdk:"variable"`
	Default     types.Map    `tfsdk:"default"`
}

func (m *CreativeTemplateResourceModel) createRequestBody() map[string]interface{} {
	body := make(map[string]interface{}, 4)
	AddInt64ValueToMap(&body, "Id", m.Id)
	AddStringValueToMap(&body, "Name", m.Name)
	return body
}

func setStateWithCreativeTemplate(s *tfsdk.State, ctx context.Context, creativeTemplate *kevelManagementClient.CreativeTemplate) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if creativeTemplate == nil {
		diags.AddError("Error", "creative template is nil")
		return diags
	}

	if creativeTemplate.Id != nil {
		diags.Append(s.SetAttribute(ctx, path.Root("id"), NewInt64PointerValueFromInt32(creativeTemplate.Id))...)
	}

	diags.Append(s.SetAttribute(ctx, path.Root("name"), types.StringValue(creativeTemplate.Name))...)

	diags.Append(s.SetAttribute(ctx, path.Root("description"), types.StringValue(creativeTemplate.Description))...)

	return diags
}

func (r *CreativeTemplateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ad_type"
}

func (r *CreativeTemplateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Kevel CreativeTemplate",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the creative template",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the creative template",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *CreativeTemplateResource) Configure(_ context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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

func (r *CreativeTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CreativeTemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.CreateCreativeTemplateWithResponse(ctx, plan.createRequestBody())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating creative template",
			"Could not create creative template, unexpected error: "+err.Error(),
		)
		return
	}

	CreativeTemplate := response.JSON200

	resp.Diagnostics.Append(setStateWithCreativeTemplate(&resp.State, ctx, CreativeTemplate)...)
}

func (r *CreativeTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CreativeTemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.GetCreativeTemplateWithResponse(ctx, int32(state.Id.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kevel CreativeTemplate",
			"Could not read creative template ID "+state.Id.String()+", unexpected error: "+err.Error(),
		)
		return
	}

	creativeTemplate := response.JSON200

	resp.Diagnostics.Append(setStateWithCreativeTemplate(&resp.State, ctx, creativeTemplate)...)
}

func (r *CreativeTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CreativeTemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.CreateCreativeTemplateWithResponse(ctx, plan.createRequestBody())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Kevel CreativeTemplate",
			"Could not update creative template ID "+plan.Id.String()+", unexpected error: "+err.Error(),
		)
		return
	}

	CreativeTemplate := response.JSON200

	resp.Diagnostics.Append(setStateWithCreativeTemplate(&resp.State, ctx, CreativeTemplate)...)
}

func (r *CreativeTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CreativeTemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := make(map[string]interface{})
	body["Updates"] = []map[string]interface{}{
		{"Path": []interface{}{}, "Op": "Delete"},
	}

	response, err := r.client.UpdateCreativeTemplateWithResponse(ctx, int32(state.Id.ValueInt64()), body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Kevel CreativeTemplate",
			"Could not delete creative template ID "+state.Id.String()+", unexpected error: "+err.Error(),
		)
		return
	}

	if response.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Error Deleting Kevel CreativeTemplate",
			"Could not delete creative template ID "+state.Id.String()+", unexpected status code: "+strconv.Itoa(response.StatusCode()),
		)
		return
	}
}

func (r *CreativeTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseUint(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing Kevel CreativeTemplate",
			"Could not import creative template ID "+req.ID+", unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
