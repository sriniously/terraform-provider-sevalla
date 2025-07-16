// Package provider contains the Terraform provider implementation for Sevalla.
package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sriniously/terraform-provider-sevalla/internal/sevallaapi"
)

var _ datasource.DataSource = &ApplicationDataSource{}

func NewApplicationDataSource() datasource.DataSource {
	return &ApplicationDataSource{}
}

type ApplicationDataSource struct {
	client *sevallaapi.Client
}

func (d *ApplicationDataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (d *ApplicationDataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches information about a Sevalla application.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Application identifier",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Application name",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Application description",
				Computed:            true,
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "Custom domain for the application",
				Computed:            true,
			},
			"repository": schema.SingleNestedAttribute{
				MarkdownDescription: "Source code repository configuration",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"url": schema.StringAttribute{
						MarkdownDescription: "Repository URL",
						Computed:            true,
					},
					"type": schema.StringAttribute{
						MarkdownDescription: "Repository type",
						Computed:            true,
					},
					"branch": schema.StringAttribute{
						MarkdownDescription: "Repository branch",
						Computed:            true,
					},
				},
			},
			"branch": schema.StringAttribute{
				MarkdownDescription: "Git branch to deploy",
				Computed:            true,
			},
			"build_command": schema.StringAttribute{
				MarkdownDescription: "Build command to run",
				Computed:            true,
			},
			"start_command": schema.StringAttribute{
				MarkdownDescription: "Start command to run",
				Computed:            true,
			},
			"environment": schema.MapAttribute{
				MarkdownDescription: "Environment variables",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"instances": schema.Int64Attribute{
				MarkdownDescription: "Number of instances",
				Computed:            true,
			},
			"memory": schema.Int64Attribute{
				MarkdownDescription: "Memory allocation in MB",
				Computed:            true,
			},
			"cpu": schema.Int64Attribute{
				MarkdownDescription: "CPU allocation in millicores",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Application status",
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

func (d *ApplicationDataSource) Configure(
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

func (d *ApplicationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ApplicationResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "reading application data source")

	app, err := sevallaapi.NewApplicationService(d.client).Get(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read application, got error: %s", err))
		return
	}

	// Use the same updateModelFromAPI function from the resource
	appResource := &ApplicationResource{}
	appResource.updateModelFromAPI(ctx, &data, app)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}