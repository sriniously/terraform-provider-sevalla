package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sriniously/terraform-provider-sevalla/internal/sevallaapi"
)

var _ datasource.DataSource = &StaticSiteDataSource{}

func NewStaticSiteDataSource() datasource.DataSource {
	return &StaticSiteDataSource{}
}

type StaticSiteDataSource struct {
	client *sevallaapi.Client
}

func (d *StaticSiteDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_static_site"
}

func (d *StaticSiteDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches information about a Sevalla static site.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Static site identifier",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Static site name",
				Computed:            true,
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "Custom domain for the static site",
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
			"build_dir": schema.StringAttribute{
				MarkdownDescription: "Build output directory",
				Computed:            true,
			},
			"build_cmd": schema.StringAttribute{
				MarkdownDescription: "Build command to run",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Static site status",
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

func (d *StaticSiteDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(SevallaProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected SevallaProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client.Client
}

func (d *StaticSiteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data StaticSiteResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "reading static site data source")

	site, err := sevallaapi.NewStaticSiteService(d.client).Get(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read static site, got error: %s", err))
		return
	}

	// Use the same updateModelFromAPI function from the resource
	siteResource := &StaticSiteResource{}
	siteResource.updateModelFromAPI(ctx, &data, site)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}