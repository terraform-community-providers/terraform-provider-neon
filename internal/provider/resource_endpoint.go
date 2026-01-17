package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &EndpointResource{}
var _ resource.ResourceWithImportState = &EndpointResource{}

func NewEndpointResource() resource.Resource {
	return &EndpointResource{}
}

type EndpointResource struct {
	client *http.Client
}

type EndpointResourceModel struct {
	Id                 types.String  `tfsdk:"id"`
	BranchId           types.String  `tfsdk:"branch_id"`
	ProjectId          types.String  `tfsdk:"project_id"`
	Type               types.String  `tfsdk:"type"`
	Host               types.String  `tfsdk:"host"`
	MinCu              types.Float64 `tfsdk:"min_cu"`
	MaxCu              types.Float64 `tfsdk:"max_cu"`
	ComputeProvisioner types.String  `tfsdk:"compute_provisioner"`
	SuspendTimeout     types.Int64   `tfsdk:"suspend_timeout"`
}

func (r *EndpointResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_endpoint"
}

func (r *EndpointResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Neon endpoint. This creates `read_only` endpoints. Please use `neon_branch` to create the `read_write` endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the endpoint.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"branch_id": schema.StringAttribute{
				MarkdownDescription: "Branch the endpoint belongs to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(idRegex(), "must be an id"),
				},
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "Project the endpoint belongs to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(idRegex(), "must be an id"),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of the endpoint.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("read_only", "read_write"),
				},
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "Host of the endpoint.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"min_cu": schema.Float64Attribute{
				MarkdownDescription: "Minimum number of compute units for the endpoint. **Default** `0.25`.",
				Optional:            true,
				Computed:            true,
				Default:             float64default.StaticFloat64(0.25),
				Validators: []validator.Float64{
					ComputeSizeValidator,
				},
			},
			"max_cu": schema.Float64Attribute{
				MarkdownDescription: "Maximum number of compute units for the endpoint. **Default** `0.25`.",
				Optional:            true,
				Computed:            true,
				Default:             float64default.StaticFloat64(0.25),
				Validators: []validator.Float64{
					ComputeSizeValidator,
				},
			},
			"compute_provisioner": schema.StringAttribute{
				MarkdownDescription: "Provisioner of the endpoint.",
				Computed:            true,
			},
			"suspend_timeout": schema.Int64Attribute{
				MarkdownDescription: "Suspend timeout of the endpoint. **Default** `0`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.Between(-1, 604800),
				},
			},
		},
	}
}

func (r *EndpointResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *EndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *EndpointResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := EndpointCreateInput{
		Endpoint: EndpointCreateInputEndpoint{
			BranchId:              data.BranchId.ValueString(),
			Type:                  "read_only",
			AutoscalingLimitMinCu: data.MinCu.ValueFloat64(),
			AutoscalingLimitMaxCu: data.MaxCu.ValueFloat64(),
			SuspendTimeoutSeconds: data.SuspendTimeout.ValueInt64(),
		},
	}

	endpoint, err := endpointCreate(r.client, data.ProjectId.ValueString(), input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create endpoint, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created a endpoint")

	data.Id = types.StringValue(endpoint.Endpoint.Id)
	data.BranchId = types.StringValue(endpoint.Endpoint.BranchId)
	data.ProjectId = types.StringValue(endpoint.Endpoint.ProjectId)
	data.Type = types.StringValue(endpoint.Endpoint.Type)
	data.Host = types.StringValue(endpoint.Endpoint.Host)
	data.MinCu = types.Float64Value(endpoint.Endpoint.AutoscalingLimitMinCu)
	data.MaxCu = types.Float64Value(endpoint.Endpoint.AutoscalingLimitMaxCu)
	data.ComputeProvisioner = types.StringValue(endpoint.Endpoint.ComputeProvisioner)
	data.SuspendTimeout = types.Int64Value(endpoint.Endpoint.SuspendTimeoutSeconds)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *EndpointResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var endpoint EndpointOutput

	err := get(r.client, fmt.Sprintf("/projects/%s/endpoints/%s", data.ProjectId.ValueString(), data.Id.ValueString()), &endpoint)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read endpoint, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "read a endpoint")

	data.Id = types.StringValue(endpoint.Endpoint.Id)
	data.BranchId = types.StringValue(endpoint.Endpoint.BranchId)
	data.ProjectId = types.StringValue(endpoint.Endpoint.ProjectId)
	data.Type = types.StringValue(endpoint.Endpoint.Type)
	data.Host = types.StringValue(endpoint.Endpoint.Host)
	data.MinCu = types.Float64Value(endpoint.Endpoint.AutoscalingLimitMinCu)
	data.MaxCu = types.Float64Value(endpoint.Endpoint.AutoscalingLimitMaxCu)
	data.ComputeProvisioner = types.StringValue(endpoint.Endpoint.ComputeProvisioner)
	data.SuspendTimeout = types.Int64Value(endpoint.Endpoint.SuspendTimeoutSeconds)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *EndpointResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := EndpointUpdateInput{
		Endpoint: EndpointUpdateInputEndpoint{
			AutoscalingLimitMinCu: data.MinCu.ValueFloat64(),
			AutoscalingLimitMaxCu: data.MaxCu.ValueFloat64(),
			SuspendTimeoutSeconds: data.SuspendTimeout.ValueInt64(),
		},
	}

	endpoint, err := endpointUpdate(r.client, data.ProjectId.ValueString(), data.Id.ValueString(), input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update endpoint, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated a endpoint")

	data.Id = types.StringValue(endpoint.Endpoint.Id)
	data.BranchId = types.StringValue(endpoint.Endpoint.BranchId)
	data.ProjectId = types.StringValue(endpoint.Endpoint.ProjectId)
	data.Type = types.StringValue(endpoint.Endpoint.Type)
	data.Host = types.StringValue(endpoint.Endpoint.Host)
	data.MinCu = types.Float64Value(endpoint.Endpoint.AutoscalingLimitMinCu)
	data.MaxCu = types.Float64Value(endpoint.Endpoint.AutoscalingLimitMaxCu)
	data.ComputeProvisioner = types.StringValue(endpoint.Endpoint.ComputeProvisioner)
	data.SuspendTimeout = types.Int64Value(endpoint.Endpoint.SuspendTimeoutSeconds)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *EndpointResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := endpointDelete(r.client, data.ProjectId.ValueString(), data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete endpoint, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted a endpoint")
}

func (r *EndpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: project_id:endpoint_id. Got: %q", req.ID),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
