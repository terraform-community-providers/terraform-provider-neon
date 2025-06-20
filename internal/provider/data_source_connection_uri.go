package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ConnectionURIDataSource{}

func NewConnectionURIDataSource() datasource.DataSource {
	return &ConnectionURIDataSource{}
}

type ConnectionURIDataSource struct {
	client *http.Client
}

type ConnectionURIDataSourceModel struct {
	Id           types.String `tfsdk:"id"`
	ProjectId    types.String `tfsdk:"project_id"`
	BranchId     types.String `tfsdk:"branch_id"`
	EndpointId   types.String `tfsdk:"endpoint_id"`
	DatabaseName types.String `tfsdk:"database_name"`
	RoleName     types.String `tfsdk:"role_name"`
	URI          types.String `tfsdk:"uri"`
	PooledURI    types.String `tfsdk:"pooled_uri"`
}

func (d *ConnectionURIDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connection_uri"
}

func (d *ConnectionURIDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves normal and pooled connection URIs for a Neon database and role.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the connection URI data source.",
				Computed:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "Project to retrieve the connection URI for.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(idRegex(), "must be an id"),
				},
			},
			"branch_id": schema.StringAttribute{
				MarkdownDescription: "Branch to retrieve the connection URI for. Defaults to the project's default branch.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(idRegex(), "must be an id"),
				},
			},
			"endpoint_id": schema.StringAttribute{
				MarkdownDescription: "Endpoint to retrieve the connection URI for. Defaults to the banch's read-write endpoint.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(idRegex(), "must be an id"),
				},
			},
			"database_name": schema.StringAttribute{
				MarkdownDescription: "Name of the database.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"role_name": schema.StringAttribute{
				MarkdownDescription: "Name of the role.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"uri": schema.StringAttribute{
				MarkdownDescription: "Normal connection URI.",
				Computed:            true,
			},
			"pooled_uri": schema.StringAttribute{
				MarkdownDescription: "Pooled connection URI.",
				Computed:            true,
			},
		},
	}
}

func (d *ConnectionURIDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*http.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *ConnectionURIDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ConnectionURIDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	uri, err := connectionURI(
		d.client,
		data.ProjectId.ValueString(),
		ConnectionURIInput{
			BranchId:     data.BranchId.ValueStringPointer(),
			EndpointId:   data.EndpointId.ValueStringPointer(),
			DatabaseName: data.DatabaseName.ValueString(),
			RoleName:     data.RoleName.ValueString(),
			Pooled:       false,
		},
	)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get connection URI, got error: %s", err))
		return
	}

	pooledURI, err := connectionURI(
		d.client,
		data.ProjectId.ValueString(),
		ConnectionURIInput{
			BranchId:     data.BranchId.ValueStringPointer(),
			EndpointId:   data.EndpointId.ValueStringPointer(),
			DatabaseName: data.DatabaseName.ValueString(),
			RoleName:     data.RoleName.ValueString(),
			Pooled:       true,
		},
	)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get pooled connection URI, got error: %s", err))
		return
	}

	data.URI = types.StringValue(uri.URI)
	data.PooledURI = types.StringValue(pooledURI.URI)
	data.Id = types.StringValue(fmt.Sprintf(
		"%s:%s:%s:%s:%s",
		data.ProjectId.ValueString(),
		data.BranchId.ValueString(),
		data.EndpointId.ValueString(),
		data.DatabaseName.ValueString(),
		data.RoleName.ValueString(),
	))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
