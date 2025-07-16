package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sriniously/terraform-provider-sevalla/internal/sevallaapi"
)

var _ datasource.DataSource = &ObjectStorageDataSource{}

func NewObjectStorageDataSource() datasource.DataSource {
	return &ObjectStorageDataSource{}
}

type ObjectStorageDataSource struct {
	client *sevallaapi.Client
}

func (d *ObjectStorageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object_storage"
}

func (d *ObjectStorageDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches information about a Sevalla object storage bucket.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Object storage identifier",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Object storage bucket name",
				Computed:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "Object storage region",
				Computed:            true,
			},
			"size": schema.Int64Attribute{
				MarkdownDescription: "Total size in bytes",
				Computed:            true,
			},
			"objects": schema.Int64Attribute{
				MarkdownDescription: "Number of objects",
				Computed:            true,
			},
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Object storage endpoint",
				Computed:            true,
			},
			"access_key": schema.StringAttribute{
				MarkdownDescription: "Object storage access key",
				Computed:            true,
			},
			"secret_key": schema.StringAttribute{
				MarkdownDescription: "Object storage secret key",
				Computed:            true,
				Sensitive:           true,
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

func (d *ObjectStorageDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ObjectStorageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ObjectStorageResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "reading object storage data source")

	bucket, err := sevallaapi.NewObjectStorageService(d.client).Get(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read object storage, got error: %s", err))
		return
	}

	// Use the same updateModelFromAPI function from the resource
	bucketResource := &ObjectStorageResource{}
	bucketResource.updateModelFromAPI(ctx, &data, bucket)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}