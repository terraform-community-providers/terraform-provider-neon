package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &BranchResource{}
var _ resource.ResourceWithImportState = &BranchResource{}

func NewBranchResource() resource.Resource {
	return &BranchResource{}
}

type BranchResource struct {
	client *http.Client
}

type BranchResourceEndpointModel struct {
	Id             types.String  `tfsdk:"id"`
	Host           types.String  `tfsdk:"host"`
	MinCu          types.Float64 `tfsdk:"min_cu"`
	MaxCu          types.Float64 `tfsdk:"max_cu"`
	Provisioner    types.String  `tfsdk:"compute_provisioner"`
	SuspendTimeout types.Int64   `tfsdk:"suspend_timeout"`
}

var endpointAttrTypes = map[string]attr.Type{
	"id":                  types.StringType,
	"host":                types.StringType,
	"min_cu":              types.Float64Type,
	"max_cu":              types.Float64Type,
	"compute_provisioner": types.StringType,
	"suspend_timeout":     types.Int64Type,
}

type BranchResourceModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	ParentId  types.String `tfsdk:"parent_id"`
	ProjectId types.String `tfsdk:"project_id"`
	Endpoint  types.Object `tfsdk:"endpoint"`
}

func BranchProvisionerCalculator() planmodifier.String {
	return branchProvisionerCalculatorModifier{}
}

type branchProvisionerCalculatorModifier struct{}

func (m branchProvisionerCalculatorModifier) Description(_ context.Context) string {
	return "This will be calculated based on compute units."
}

func (m branchProvisionerCalculatorModifier) MarkdownDescription(_ context.Context) string {
	return "This will be calculated based on compute units."
}

func (m branchProvisionerCalculatorModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	var data *BranchResourceModel
	var endpointData *BranchResourceEndpointModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(data.Endpoint.As(ctx, &endpointData, basetypes.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return
	}

	if endpointData.MinCu == endpointData.MaxCu {
		resp.PlanValue = types.StringValue("k8s-pod")
	} else {
		resp.PlanValue = types.StringValue("k8s-neonvm")
	}
}

func (r *BranchResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_branch"
}

func (r *BranchResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Neon branch. Please use `neon_project` to create primary branch.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the branch.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the branch.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"parent_id": schema.StringAttribute{
				MarkdownDescription: "ID of the parent branch. Defaults to the primary branch.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(idRegex(), "must be an id"),
				},
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "Project the branch belongs to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(idRegex(), "must be an id"),
				},
			},
			"endpoint": schema.SingleNestedAttribute{
				MarkdownDescription: "Read-write compute endpoint settings of the branch.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "Identifier of the endpoint.",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
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
							float64validator.OneOf(0.25, 0.5, 1, 2, 3, 4, 5, 6, 7),
						},
					},
					"max_cu": schema.Float64Attribute{
						MarkdownDescription: "Maximum number of compute units for the endpoint. **Default** `0.25`.",
						Optional:            true,
						Computed:            true,
						Default:             float64default.StaticFloat64(0.25),
						Validators: []validator.Float64{
							float64validator.OneOf(0.25, 0.5, 1, 2, 3, 4, 5, 6, 7),
						},
					},
					"compute_provisioner": schema.StringAttribute{
						MarkdownDescription: "Provisioner of the endpoint.",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							BranchProvisionerCalculator(),
						},
					},
					"suspend_timeout": schema.Int64Attribute{
						MarkdownDescription: "Suspend timeout of the endpoint. **Default** `300`.",
						Optional:            true,
						Computed:            true,
						Default:             int64default.StaticInt64(300),
						Validators: []validator.Int64{
							int64validator.Between(-1, 604800),
						},
					},
				},
			},
		},
	}
}

func (r *BranchResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *BranchResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *BranchResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := BranchCreateInput{
		Branch: BranchCreateInputBranch{
			Name: data.Name.ValueString(),
		},
	}

	if !data.ParentId.IsUnknown() {
		value := data.ParentId.ValueString()
		input.Branch.ParentId = value
	}

	branch, err := branchCreate(r.client, data.ProjectId.ValueString(), input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create branch, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created a branch")

	data.Id = types.StringValue(branch.Branch.Id)
	data.Name = types.StringValue(branch.Branch.Name)
	data.ProjectId = types.StringValue(branch.Branch.ProjectId)

	if branch.Branch.ParentId != nil {
		data.ParentId = types.StringValue(*branch.Branch.ParentId)
	} else {
		data.ParentId = types.StringNull()
	}

	// Create endpoint if endpoint is set
	if !data.Endpoint.IsNull() {
		var endpointData *BranchResourceEndpointModel

		resp.Diagnostics.Append(data.Endpoint.As(ctx, &endpointData, basetypes.ObjectAsOptions{})...)

		if resp.Diagnostics.HasError() {
			return
		}

		input := EndpointCreateInput{
			Endpoint: EndpointCreateInputEndpoint{
				BranchId:              branch.Branch.Id,
				Type:                  "read_write",
				AutoscalingLimitMinCu: endpointData.MinCu.ValueFloat64(),
				AutoscalingLimitMaxCu: endpointData.MaxCu.ValueFloat64(),
				Provisioner:           endpointData.Provisioner.ValueString(),
				SuspendTimeoutSeconds: endpointData.SuspendTimeout.ValueInt64(),
			},
		}

		endpoint, err := endpointCreate(r.client, data.ProjectId.ValueString(), input)

		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create endpoint of the branch, got error: %s", err))
			return
		}

		tflog.Trace(ctx, "created an endpoint")

		data.Endpoint = types.ObjectValueMust(
			endpointAttrTypes,
			map[string]attr.Value{
				"id":                  types.StringValue(endpoint.Endpoint.Id),
				"host":                types.StringValue(endpoint.Endpoint.Host),
				"min_cu":              types.Float64Value(endpoint.Endpoint.AutoscalingLimitMinCu),
				"max_cu":              types.Float64Value(endpoint.Endpoint.AutoscalingLimitMaxCu),
				"compute_provisioner": types.StringValue(endpoint.Endpoint.Provisioner),
				"suspend_timeout":     types.Int64Value(endpoint.Endpoint.SuspendTimeoutSeconds),
			},
		)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BranchResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *BranchResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var branch BranchOutput

	err := get(r.client, fmt.Sprintf("/projects/%s/branches/%s", data.ProjectId.ValueString(), data.Id.ValueString()), &branch)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read branch, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "read a branch")

	endpoint, err := branchEndpoint(r.client, branch.Branch.ProjectId, branch.Branch.Id, false)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read endpoint of the branch, got error: %s", err))
		return
	}

	data.Id = types.StringValue(branch.Branch.Id)
	data.Name = types.StringValue(branch.Branch.Name)
	data.ProjectId = types.StringValue(branch.Branch.ProjectId)

	if branch.Branch.ParentId != nil {
		data.ParentId = types.StringValue(*branch.Branch.ParentId)
	} else {
		data.ParentId = types.StringNull()
	}

	if len(endpoint.Id) > 0 {
		data.Endpoint = types.ObjectValueMust(
			endpointAttrTypes,
			map[string]attr.Value{
				"id":                  types.StringValue(endpoint.Id),
				"host":                types.StringValue(endpoint.Host),
				"min_cu":              types.Float64Value(endpoint.AutoscalingLimitMinCu),
				"max_cu":              types.Float64Value(endpoint.AutoscalingLimitMaxCu),
				"compute_provisioner": types.StringValue(endpoint.Provisioner),
				"suspend_timeout":     types.Int64Value(endpoint.SuspendTimeoutSeconds),
			},
		)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BranchResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *BranchResourceModel
	var state *BranchResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var branch = Branch{
		Id:        state.Id.ValueString(),
		Name:      state.Name.ValueString(),
		ParentId:  state.ParentId.ValueStringPointer(),
		ProjectId: state.ProjectId.ValueString(),
	}

	// Need to do this check because we can't update the branch with the same name
	if data.Name.ValueString() != state.Name.ValueString() {
		input := BranchUpdateInput{
			Branch: BranchUpdateInputBranch{
				Name: data.Name.ValueString(),
			},
		}

		branchOutput, err := branchUpdate(r.client, data.ProjectId.ValueString(), data.Id.ValueString(), input)

		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update branch, got error: %s", err))
			return
		}

		tflog.Trace(ctx, "updated a branch")

		branch = branchOutput.Branch
	}

	data.Id = types.StringValue(branch.Id)
	data.Name = types.StringValue(branch.Name)
	data.ProjectId = types.StringValue(branch.ProjectId)

	if branch.ParentId != nil {
		data.ParentId = types.StringValue(*branch.ParentId)
	} else {
		data.ParentId = types.StringNull()
	}

	var endpoint Endpoint

	// Create endpoint if endpoint is set in plan and endpoint is not set in state
	if !data.Endpoint.IsNull() && state.Endpoint.IsNull() {
		var endpointData *BranchResourceEndpointModel

		resp.Diagnostics.Append(data.Endpoint.As(ctx, &endpointData, basetypes.ObjectAsOptions{})...)

		if resp.Diagnostics.HasError() {
			return
		}

		input := EndpointCreateInput{
			Endpoint: EndpointCreateInputEndpoint{
				BranchId:              branch.Id,
				Type:                  "read_write",
				AutoscalingLimitMinCu: endpointData.MinCu.ValueFloat64(),
				AutoscalingLimitMaxCu: endpointData.MaxCu.ValueFloat64(),
				Provisioner:           endpointData.Provisioner.ValueString(),
				SuspendTimeoutSeconds: endpointData.SuspendTimeout.ValueInt64(),
			},
		}

		endpointOutput, err := endpointCreate(r.client, data.ProjectId.ValueString(), input)

		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create endpoint of the branch, got error: %s", err))
			return
		}

		tflog.Trace(ctx, "created an endpoint")

		endpoint = endpointOutput.Endpoint
	}

	// Update endpoint if endpoint is set in both plan and state
	if !data.Endpoint.IsNull() && !state.Endpoint.IsNull() {
		var endpointData *BranchResourceEndpointModel

		resp.Diagnostics.Append(data.Endpoint.As(ctx, &endpointData, basetypes.ObjectAsOptions{})...)

		if resp.Diagnostics.HasError() {
			return
		}

		input := EndpointUpdateInput{
			Endpoint: EndpointUpdateInputEndpoint{
				AutoscalingLimitMinCu: endpointData.MinCu.ValueFloat64(),
				AutoscalingLimitMaxCu: endpointData.MaxCu.ValueFloat64(),
				Provisioner:           endpointData.Provisioner.ValueString(),
				SuspendTimeoutSeconds: endpointData.SuspendTimeout.ValueInt64(),
			},
		}

		endpointOuput, err := endpointUpdate(r.client, data.ProjectId.ValueString(), endpointData.Id.ValueString(), input)

		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update endpoint of the branch, got error: %s", err))
			return
		}

		tflog.Trace(ctx, "updated an endpoint")

		endpoint = endpointOuput.Endpoint
	}

	// Delete endpoint if endpoint is set in state and endpoint is not set in plan
	if data.Endpoint.IsNull() && !state.Endpoint.IsNull() {
		var endpointState *BranchResourceEndpointModel

		resp.Diagnostics.Append(state.Endpoint.As(ctx, &endpointState, basetypes.ObjectAsOptions{})...)

		if resp.Diagnostics.HasError() {
			return
		}

		err := endpointDelete(r.client, state.ProjectId.ValueString(), endpointState.Id.ValueString())

		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete endpoint of the branch, got error: %s", err))
			return
		}

		tflog.Trace(ctx, "deleted an endpoint")
	}

	if len(endpoint.Id) > 0 {
		data.Endpoint = types.ObjectValueMust(
			endpointAttrTypes,
			map[string]attr.Value{
				"id":                  types.StringValue(endpoint.Id),
				"host":                types.StringValue(endpoint.Host),
				"min_cu":              types.Float64Value(endpoint.AutoscalingLimitMinCu),
				"max_cu":              types.Float64Value(endpoint.AutoscalingLimitMaxCu),
				"compute_provisioner": types.StringValue(endpoint.Provisioner),
				"suspend_timeout":     types.Int64Value(endpoint.SuspendTimeoutSeconds),
			},
		)
	} else {
		data.Endpoint = types.ObjectNull(endpointAttrTypes)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BranchResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *BranchResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := branchDelete(r.client, data.ProjectId.ValueString(), data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete branch, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted a branch")
}

func (r *BranchResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: project_id:branch_id. Got: %q", req.ID),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
