package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	adzerk "github.com/cysp/adzerk-management-sdk-go"
)

var (
	_ resource.Resource                = &creativeTemplateResource{}
	_ resource.ResourceWithConfigure   = &creativeTemplateResource{}
	_ resource.ResourceWithImportState = &creativeTemplateResource{}
)

func NewCreativeTemplateResource() resource.Resource {
	return &creativeTemplateResource{}
}

type creativeTemplateResource struct {
	client *adzerk.ClientWithResponses
}

func (r *creativeTemplateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_creative_template"
}

func (r *creativeTemplateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Kevel CreativeTemplate",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the creative template",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the creative template",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the creative template",
				Optional:    true,
			},
			"fields": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("Array", "ExternalFile", "File", "Object", "String"),
							},
						},
						"name": schema.StringAttribute{
							Required: true,
						},
						"description": schema.StringAttribute{
							Optional: true,
						},
						"variable": schema.StringAttribute{
							Required: true,
						},
						"required": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
						},
						"hidden": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
						},
						"ad_query": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
						},
						"default": schema.SingleNestedAttribute{
							Attributes: map[string]schema.Attribute{
								"string": schema.StringAttribute{
									Optional: true,
								},
								"array": schema.ListAttribute{
									ElementType: types.StringType,
									Optional:    true,
								},
							},
							Optional: true,
						},
					},
				},
				Optional: true,
			},
			"contents": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("CSS", "HTML", "JavaScript", "JavaScriptExternal", "Raw"),
							},
						},
						"body": schema.StringAttribute{
							Optional: true,
						},
					},
				},
				Optional: true,
			},
		},
	}
}

func (r *creativeTemplateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *creativeTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan creativeTemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, bodyDiags := plan.createRequestBody(ctx)
	resp.Diagnostics.Append(bodyDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.CreateCreativeTemplateWithResponse(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating creative_template",
			"Could not create creative_template, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(setStateWithCreativeTemplate(&resp.State, ctx, response.JSON200)...)
}

func (r *creativeTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state creativeTemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.GetCreativeTemplateWithResponse(ctx, int32(state.Id.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kevel CreativeTemplate",
			"Could not read creative_template ID "+state.Id.String()+", unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(setStateWithCreativeTemplate(&resp.State, ctx, response.JSON200)...)
}

func (r *creativeTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan creativeTemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.UpdateCreativeTemplateWithResponse(ctx, int32(plan.Id.ValueInt64()), plan.updateRequestBody())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Kevel CreativeTemplate",
			"Could not update creative_template ID "+plan.Id.String()+", unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(setStateWithCreativeTemplate(&resp.State, ctx, response.JSON200)...)
}

func (r *creativeTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state creativeTemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.UpdateCreativeTemplateWithResponse(ctx, int32(state.Id.ValueInt64()), state.deleteRequestBody())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Kevel CreativeTemplate",
			"Could not delete creative_template ID "+state.Id.String()+", unexpected error: "+err.Error(),
		)
		return
	}

	if response.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"Error Deleting Kevel CreativeTemplate",
			"Could not delete creative_template ID "+state.Id.String()+", unexpected status code: "+strconv.Itoa(response.StatusCode()),
		)
		return
	}
}

func (r *creativeTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ImportStatePassthroughInt64ID(ctx, path.Root("id"), req, resp)
}
