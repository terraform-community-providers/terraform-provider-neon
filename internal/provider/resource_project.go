package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/exp/slices"
)

var _ resource.Resource = &ProjectResource{}
var _ resource.ResourceWithImportState = &ProjectResource{}

func NewProjectResource() resource.Resource {
	return &ProjectResource{}
}

type ProjectResource struct {
	client *http.Client
}

type ProjectResourceBranchEndpointModel struct {
	Id                 types.String  `tfsdk:"id"`
	Host               types.String  `tfsdk:"host"`
	MinCu              types.Float64 `tfsdk:"min_cu"`
	MaxCu              types.Float64 `tfsdk:"max_cu"`
	ComputeProvisioner types.String  `tfsdk:"compute_provisioner"`
	SuspendTimeout     types.Int64   `tfsdk:"suspend_timeout"`
}

var branchEndpointAttrTypes = map[string]attr.Type{
	"id":                  types.StringType,
	"host":                types.StringType,
	"min_cu":              types.Float64Type,
	"max_cu":              types.Float64Type,
	"compute_provisioner": types.StringType,
	"suspend_timeout":     types.Int64Type,
}

type ProjectResourceBranchModel struct {
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Endpoint types.Object `tfsdk:"endpoint"`
}

var branchAttrTypes = map[string]attr.Type{
	"id":   types.StringType,
	"name": types.StringType,
	"endpoint": types.ObjectType{
		AttrTypes: branchEndpointAttrTypes,
	},
}

type ProjectResourceModel struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	PlatformId types.String `tfsdk:"platform_id"`
	RegionId   types.String `tfsdk:"region_id"`
	PgVersion  types.Int64  `tfsdk:"pg_version"`
	Branch     types.Object `tfsdk:"branch"`
}

func ProjectProvisionerCalculator() planmodifier.String {
	return projectProvisionerCalculatorModifier{}
}

type projectProvisionerCalculatorModifier struct{}

func (m projectProvisionerCalculatorModifier) Description(_ context.Context) string {
	return "This will be calculated based on compute units."
}

func (m projectProvisionerCalculatorModifier) MarkdownDescription(_ context.Context) string {
	return "This will be calculated based on compute units."
}

func (m projectProvisionerCalculatorModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	var data *ProjectResourceModel
	var branchData *ProjectResourceBranchModel
	var branchEndpointData *ProjectResourceBranchEndpointModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(data.Branch.As(ctx, &branchData, basetypes.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(branchData.Endpoint.As(ctx, &branchEndpointData, basetypes.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return
	}

	if branchEndpointData.MinCu == branchEndpointData.MaxCu {
		resp.PlanValue = types.StringValue("k8s-pod")
	} else {
		resp.PlanValue = types.StringValue("k8s-neonvm")
	}
}

func (r *ProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *ProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Neon project.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the project.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the project.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
					stringvalidator.UTF8LengthAtMost(64),
				},
			},
			"platform_id": schema.StringAttribute{
				MarkdownDescription: "Platform of the project.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region_id": schema.StringAttribute{
				MarkdownDescription: "Region of the project.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"pg_version": schema.Int64Attribute{
				MarkdownDescription: "PostgreSQL version of the project.",
				Computed:            true,
				Optional:            true,
				Default:             int64default.StaticInt64(15),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.OneOf(14, 15),
				},
			},
			"branch": schema.SingleNestedAttribute{
				MarkdownDescription: "Primary branch settings of the project.",
				Optional:            true,
				Computed:            true,
				Default: objectdefault.StaticValue(
					types.ObjectValueMust(
						branchAttrTypes,
						map[string]attr.Value{
							"id":   types.StringUnknown(),
							"name": types.StringValue("main"),
							"endpoint": types.ObjectValueMust(
								branchEndpointAttrTypes,
								map[string]attr.Value{
									"id":                  types.StringUnknown(),
									"host":                types.StringUnknown(),
									"min_cu":              types.Float64Value(0.25),
									"max_cu":              types.Float64Value(0.25),
									"compute_provisioner": types.StringUnknown(),
									"suspend_timeout":     types.Int64Value(300),
								},
							),
						},
					),
				),
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "Identifier of the branch.",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "Name of the branch.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("main"),
						Validators: []validator.String{
							stringvalidator.UTF8LengthAtLeast(1),
						},
					},
					"endpoint": schema.SingleNestedAttribute{
						MarkdownDescription: "Read-write compute endpoint settings of the branch.",
						Optional:            true,
						Computed:            true,
						Default: objectdefault.StaticValue(
							types.ObjectValueMust(
								branchEndpointAttrTypes,
								map[string]attr.Value{
									"id":                  types.StringUnknown(),
									"host":                types.StringUnknown(),
									"min_cu":              types.Float64Value(0.25),
									"max_cu":              types.Float64Value(0.25),
									"compute_provisioner": types.StringUnknown(),
									"suspend_timeout":     types.Int64Value(300),
								},
							),
						),
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
									ProjectProvisionerCalculator(),
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
			},
		},
	}
}

func (r *ProjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ProjectResourceModel
	var branchData *ProjectResourceBranchModel
	var branchEndpointData *ProjectResourceBranchEndpointModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := ProjectCreateInput{
		Project: ProjectCreateInputProject{
			Name:           data.Name.ValueString(),
			RegionId:       data.RegionId.ValueString(),
			PgVersion:      data.PgVersion.ValueInt64(),
			StorePasswords: true,
		},
	}

	resp.Diagnostics.Append(data.Branch.As(ctx, &branchData, basetypes.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return
	}

	input.Project.Branch = ProjectCreateInputProjectBranch{
		Name: branchData.Name.ValueString(),
	}

	resp.Diagnostics.Append(branchData.Endpoint.As(ctx, &branchEndpointData, basetypes.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return
	}

	input.Project.ComputeProvisioner = branchEndpointData.ComputeProvisioner.ValueString()

	input.Project.DefaultEndpointSettings = ProjectCreateInputProjectDefaultEndpointSettings{
		AutoscalingLimitMinCu: branchEndpointData.MinCu.ValueFloat64(),
		AutoscalingLimitMaxCu: branchEndpointData.MaxCu.ValueFloat64(),
		SuspendTimeoutSeconds: branchEndpointData.SuspendTimeout.ValueInt64(),
	}

	var project ProjectCreateOutput

	err := call(r.client, http.MethodPost, "/projects", input, &project)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create project, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created a project")

	// Delete the default database.
	err = databaseDelete(r.client, project.Project.Id, project.Branch.Id, project.Databases[0].Name)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete default database, got error: %s", err))
		return
	}

	// Delete the default role.
	err = roleDelete(r.client, project.Project.Id, project.Branch.Id, project.Roles[0].Name)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete default role, got error: %s", err))
		return
	}

	data.Id = types.StringValue(project.Project.Id)
	data.Name = types.StringValue(project.Project.Name)
	data.PlatformId = types.StringValue(project.Project.PlatformId)
	data.RegionId = types.StringValue(project.Project.RegionId)
	data.PgVersion = types.Int64Value(project.Project.PgVersion)

	data.Branch = types.ObjectValueMust(
		branchAttrTypes,
		map[string]attr.Value{
			"id":   types.StringValue(project.Branch.Id),
			"name": types.StringValue(project.Branch.Name),
			"endpoint": types.ObjectValueMust(
				branchEndpointAttrTypes,
				map[string]attr.Value{
					"id":                  types.StringValue(project.Endpoints[0].Id),
					"host":                types.StringValue(project.Endpoints[0].Host),
					"min_cu":              types.Float64Value(project.Endpoints[0].AutoscalingLimitMinCu),
					"max_cu":              types.Float64Value(project.Endpoints[0].AutoscalingLimitMaxCu),
					"compute_provisioner": types.StringValue(project.Endpoints[0].ComputeProvisioner),
					"suspend_timeout":     types.Int64Value(project.Endpoints[0].SuspendTimeoutSeconds),
				},
			),
		},
	)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ProjectResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var project ProjectOutput

	err := get(r.client, fmt.Sprintf("/projects/%s", data.Id.ValueString()), &project)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read project, got error: %s", err))
		return
	}

	// Get the primary branch for the project
	branch, err := readPrimaryBranch(r.client, project.Project.Id)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read primary branch of the project, got error: %s", err))
		return
	}

	// Get the endpoint for the primary branch
	endpoint, err := branchEndpoint(r.client, project.Project.Id, branch.Id, true)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read endpoint of the primary branch, got error: %s", err))
		return
	}

	data.Id = types.StringValue(project.Project.Id)
	data.Name = types.StringValue(project.Project.Name)
	data.PlatformId = types.StringValue(project.Project.PlatformId)
	data.RegionId = types.StringValue(project.Project.RegionId)
	data.PgVersion = types.Int64Value(project.Project.PgVersion)

	data.Branch = types.ObjectValueMust(
		branchAttrTypes,
		map[string]attr.Value{
			"id":   types.StringValue(branch.Id),
			"name": types.StringValue(branch.Name),
			"endpoint": types.ObjectValueMust(
				branchEndpointAttrTypes,
				map[string]attr.Value{
					"id":                  types.StringValue(endpoint.Id),
					"host":                types.StringValue(endpoint.Host),
					"min_cu":              types.Float64Value(endpoint.AutoscalingLimitMinCu),
					"max_cu":              types.Float64Value(endpoint.AutoscalingLimitMaxCu),
					"compute_provisioner": types.StringValue(endpoint.ComputeProvisioner),
					"suspend_timeout":     types.Int64Value(endpoint.SuspendTimeoutSeconds),
				},
			),
		},
	)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ProjectResourceModel
	var branchData *ProjectResourceBranchModel
	var branchEndpointData *ProjectResourceBranchEndpointModel

	var state *ProjectResourceModel
	var branchState *ProjectResourceBranchModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := ProjectUpdateInput{
		Project: ProjectUpdateInputProject{
			Name: data.Name.ValueString(),
		},
	}

	var project ProjectOutput

	err := call(r.client, http.MethodPatch, fmt.Sprintf("/projects/%s", data.Id.ValueString()), input, &project)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update project, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated a project")

	resp.Diagnostics.Append(data.Branch.As(ctx, &branchData, basetypes.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(state.Branch.As(ctx, &branchState, basetypes.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return
	}

	branch := Branch{
		Id:   branchState.Id.ValueString(),
		Name: branchState.Name.ValueString(),
	}

	// Need to do this check because we can't update the branch with the same name
	if branchData.Name.ValueString() != branchState.Name.ValueString() {
		branchInput := BranchUpdateInput{
			Branch: BranchUpdateInputBranch{
				Name: branchData.Name.ValueString(),
			},
		}

		branchOutput, err := branchUpdate(r.client, data.Id.ValueString(), branchData.Id.ValueString(), branchInput)

		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update branch, got error: %s", err))
			return
		}

		tflog.Trace(ctx, "updated a branch")

		branch = branchOutput.Branch
	}

	resp.Diagnostics.Append(branchData.Endpoint.As(ctx, &branchEndpointData, basetypes.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return
	}

	endpointInput := EndpointUpdateInput{
		Endpoint: EndpointUpdateInputEndpoint{
			AutoscalingLimitMinCu: branchEndpointData.MinCu.ValueFloat64(),
			AutoscalingLimitMaxCu: branchEndpointData.MaxCu.ValueFloat64(),
			ComputeProvisioner:    branchEndpointData.ComputeProvisioner.ValueString(),
			SuspendTimeoutSeconds: branchEndpointData.SuspendTimeout.ValueInt64(),
		},
	}

	endpoint, err := endpointUpdate(r.client, data.Id.ValueString(), branchEndpointData.Id.ValueString(), endpointInput)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update endpoint, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated an endpoint")

	data.Id = types.StringValue(project.Project.Id)
	data.Name = types.StringValue(project.Project.Name)
	data.PlatformId = types.StringValue(project.Project.PlatformId)
	data.RegionId = types.StringValue(project.Project.RegionId)
	data.PgVersion = types.Int64Value(project.Project.PgVersion)

	data.Branch = types.ObjectValueMust(
		branchAttrTypes,
		map[string]attr.Value{
			"id":   types.StringValue(branch.Id),
			"name": types.StringValue(branch.Name),
			"endpoint": types.ObjectValueMust(
				branchEndpointAttrTypes,
				map[string]attr.Value{
					"id":                  types.StringValue(endpoint.Endpoint.Id),
					"host":                types.StringValue(endpoint.Endpoint.Host),
					"min_cu":              types.Float64Value(endpoint.Endpoint.AutoscalingLimitMinCu),
					"max_cu":              types.Float64Value(endpoint.Endpoint.AutoscalingLimitMaxCu),
					"compute_provisioner": types.StringValue(endpoint.Endpoint.ComputeProvisioner),
					"suspend_timeout":     types.Int64Value(endpoint.Endpoint.SuspendTimeoutSeconds),
				},
			),
		},
	)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ProjectResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := delete(r.client, fmt.Sprintf("/projects/%s", data.Id.ValueString()))

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete project, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted a project")
}

func (r *ProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func readPrimaryBranch(client *http.Client, projectId string) (Branch, error) {
	var branch Branch

	// Read all branches
	branches, err := branchList(client, projectId)

	if err != nil {
		return branch, err
	}

	// Get the primary branch
	branchIdx := slices.IndexFunc(branches.Branches, func(branch Branch) bool {
		return branch.Primary
	})

	return branches.Branches[branchIdx], nil
}
