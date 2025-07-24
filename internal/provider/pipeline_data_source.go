package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sriniously/terraform-provider-sevalla/internal/sevallaapi"
)

var _ datasource.DataSource = &PipelineDataSource{}

func NewPipelineDataSource() datasource.DataSource {
	return &PipelineDataSource{}
}

type PipelineDataSource struct {
	client *sevallaapi.Client
}

func (d *PipelineDataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_pipeline"
}

func (d *PipelineDataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches information about a Sevalla deployment pipeline.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Pipeline identifier",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Pipeline name",
				Computed:            true,
			},
			"app_id": schema.StringAttribute{
				MarkdownDescription: "Associated application ID",
				Computed:            true,
			},
			"branch": schema.StringAttribute{
				MarkdownDescription: "Git branch for the pipeline",
				Computed:            true,
			},
			"auto_deploy": schema.BoolAttribute{
				MarkdownDescription: "Whether automatic deployment is enabled",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Creation timestamp",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "Last update timestamp",
				Computed:            true,
			},
		},
	}
}

func (d *PipelineDataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(SevallaProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected SevallaProviderData, got: %T. "+
				"Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client.Client
}

func (d *PipelineDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PipelineResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "reading pipeline data source")

	pipeline, err := sevallaapi.NewPipelineService(d.client).Get(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read pipeline, got error: %s", err))
		return
	}

	// Use the same updateModelFromAPI function from the resource
	pipelineResource := &PipelineResource{}
	pipelineResource.updateModelFromAPI(ctx, &data, pipeline)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
