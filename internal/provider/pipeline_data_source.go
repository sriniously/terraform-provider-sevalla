package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sriniously/terraform-provider-sevalla/internal/sevallaapi"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &PipelineDataSource{}

func NewPipelineDataSource() datasource.DataSource {
	return &PipelineDataSource{}
}

// PipelineDataSource defines the data source implementation.
type PipelineDataSource struct {
	client *sevallaapi.Client
}

// PipelineDataSourceModel describes the data source data model.
type PipelineDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	AppID      types.String `tfsdk:"app_id"`
	Branch     types.String `tfsdk:"branch"`
	AutoDeploy types.Bool   `tfsdk:"auto_deploy"`
	CreatedAt  types.String `tfsdk:"created_at"`
	UpdatedAt  types.String `tfsdk:"updated_at"`
}

func (d *PipelineDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pipeline"
}

func (d *PipelineDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for fetching information about a Sevalla deployment pipeline.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier of the pipeline.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the pipeline.",
			},
			"app_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the application this pipeline deploys.",
			},
			"branch": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The git branch to deploy from.",
			},
			"auto_deploy": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether to automatically deploy when changes are pushed to the branch.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The timestamp when the pipeline was created.",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The timestamp when the pipeline was last updated.",
			},
		},
	}
}

func (d *PipelineDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(SevallaProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected SevallaProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = data.Client
}

func (d *PipelineDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PipelineDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get pipeline from API
	pipeline, err := d.client.GetPipeline(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read pipeline, got error: %s", err))
		return
	}

	// Map response back to schema
	data.ID = types.StringValue(pipeline.ID)
	data.Name = types.StringValue(pipeline.DisplayName)
	// Set other values from API response
	data.AppID = types.StringValue("")      // Set from API response when available
	data.Branch = types.StringValue("main") // Set from API response when available
	data.AutoDeploy = types.BoolValue(true) // Set from API response when available
	data.CreatedAt = types.StringValue("")  // Set from API response when available
	data.UpdatedAt = types.StringValue("")  // Set from API response when available

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
