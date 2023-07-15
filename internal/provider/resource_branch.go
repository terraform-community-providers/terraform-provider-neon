package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

type BranchResourceModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	ParentId  types.String `tfsdk:"parent_id"`
	ProjectId types.String `tfsdk:"project_id"`
}

func (r *BranchResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_branch"
}

func (r *BranchResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Neon branch.",
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

	data.Id = types.StringValue(branch.Branch.Id)
	data.Name = types.StringValue(branch.Branch.Name)
	data.ProjectId = types.StringValue(branch.Branch.ProjectId)

	if branch.Branch.ParentId != nil {
		data.ParentId = types.StringValue(*branch.Branch.ParentId)
	} else {
		data.ParentId = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BranchResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *BranchResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := BranchUpdateInput{
		Branch: BranchUpdateInputBranch{
			Name: data.Name.ValueString(),
		},
	}

	branch, err := branchUpdate(r.client, data.ProjectId.ValueString(), data.Id.ValueString(), input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update branch, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated a branch")

	data.Id = types.StringValue(branch.Branch.Id)
	data.Name = types.StringValue(branch.Branch.Name)
	data.ProjectId = types.StringValue(branch.Branch.ProjectId)

	if branch.Branch.ParentId != nil {
		data.ParentId = types.StringValue(*branch.Branch.ParentId)
	} else {
		data.ParentId = types.StringNull()
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
