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

var _ datasource.DataSource = &DatabaseDataSource{}

func NewDatabaseDataSource() datasource.DataSource {
	return &DatabaseDataSource{}
}

type DatabaseDataSource struct {
	client *sevallaapi.Client
}

func (d *DatabaseDataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_database"
}

func (d *DatabaseDataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches information about a Sevalla database.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Database identifier",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Database name",
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Database display name",
				Computed:            true,
			},
			"company_id": schema.StringAttribute{
				MarkdownDescription: "Company ID",
				Computed:            true,
			},
			"location": schema.StringAttribute{
				MarkdownDescription: "Database location",
				Computed:            true,
			},
			"resource_type": schema.StringAttribute{
				MarkdownDescription: "Resource type",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Database type",
				Computed:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "Database version",
				Computed:            true,
			},
			"db_name": schema.StringAttribute{
				MarkdownDescription: "Database name",
				Computed:            true,
			},
			"db_password": schema.StringAttribute{
				MarkdownDescription: "Database password",
				Computed:            true,
				Sensitive:           true,
			},
			"db_user": schema.StringAttribute{
				MarkdownDescription: "Database username",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Database status",
				Computed:            true,
			},
			"internal_hostname": schema.StringAttribute{
				MarkdownDescription: "Internal hostname",
				Computed:            true,
			},
			"internal_port": schema.StringAttribute{
				MarkdownDescription: "Internal port",
				Computed:            true,
			},
			"external_hostname": schema.StringAttribute{
				MarkdownDescription: "External hostname",
				Computed:            true,
			},
			"external_port": schema.StringAttribute{
				MarkdownDescription: "External port",
				Computed:            true,
			},
		},
	}
}

func (d *DatabaseDataSource) Configure(
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

func (d *DatabaseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DatabaseResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "reading database data source")

	db, err := d.client.Databases.Get(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read database, got error: %s", err))
		return
	}

	// Map API response to model
	data.ID = types.StringValue(db.Database.ID)
	data.Name = types.StringValue(db.Database.Name)
	data.DisplayName = types.StringValue(db.Database.DisplayName)
	data.CompanyID = types.StringValue("") // Not available in API response
	data.Location = types.StringValue(db.Database.Cluster.Location)
	data.ResourceType = types.StringValue(db.Database.ResourceTypeName)
	data.Type = types.StringValue(db.Database.Type)
	data.Version = types.StringValue(db.Database.Version)
	data.DBName = types.StringValue(db.Database.Data.DBName)
	data.DBPassword = types.StringValue(db.Database.Data.DBPassword)
	if db.Database.Data.DBUser != nil {
		data.DBUser = types.StringValue(*db.Database.Data.DBUser)
	}
	data.Status = types.StringValue(db.Database.Status)

	if db.Database.InternalHostname != nil {
		data.InternalHostname = types.StringValue(*db.Database.InternalHostname)
	}
	if db.Database.InternalPort != nil {
		data.InternalPort = types.StringValue(*db.Database.InternalPort)
	}
	if db.Database.ExternalHostname != nil {
		data.ExternalHostname = types.StringValue(*db.Database.ExternalHostname)
	}
	if db.Database.ExternalPort != nil {
		data.ExternalPort = types.StringValue(*db.Database.ExternalPort)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
