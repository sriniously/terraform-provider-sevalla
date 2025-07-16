package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sriniously/terraform-provider-sevalla/internal/sevallaapi"
)

var _ resource.Resource = &DatabaseResource{}
var _ resource.ResourceWithImportState = &DatabaseResource{}

func NewDatabaseResource() resource.Resource {
	return &DatabaseResource{}
}

type DatabaseResource struct {
	client *sevallaapi.Client
}

type DatabaseResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Type      types.String `tfsdk:"type"`
	Version   types.String `tfsdk:"version"`
	Size      types.String `tfsdk:"size"`
	Host      types.String `tfsdk:"host"`
	Port      types.Int64  `tfsdk:"port"`
	Username  types.String `tfsdk:"username"`
	Password  types.String `tfsdk:"password"`
	Status    types.String `tfsdk:"status"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

func (r *DatabaseResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_database"
}

func (r *DatabaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Sevalla database.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Database identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Database name",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Database type (postgresql, mysql, mariadb, redis)",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "Database version",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"size": schema.StringAttribute{
				MarkdownDescription: "Database size/plan",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Database password",
				Optional:            true,
				Sensitive:           true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "Database host",
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Database port",
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Database username",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Database status",
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

func (r *DatabaseResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(SevallaProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected SevallaProviderData, got: %T. "+
				"Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client.Client
}

func (r *DatabaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DatabaseResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := sevallaapi.CreateDatabaseRequest{
		Name: data.Name.ValueString(),
		Type: data.Type.ValueString(),
	}

	if !data.Version.IsNull() {
		createReq.Version = data.Version.ValueString()
	}

	if !data.Size.IsNull() {
		createReq.Size = data.Size.ValueString()
	}

	if !data.Password.IsNull() {
		createReq.Password = data.Password.ValueString()
	}

	tflog.Trace(ctx, "creating database")

	db, err := sevallaapi.NewDatabaseService(r.client).Create(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create database, got error: %s", err))
		return
	}

	r.updateModelFromAPI(ctx, &data, db)

	tflog.Trace(ctx, "created database")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DatabaseResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	db, err := sevallaapi.NewDatabaseService(r.client).Get(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read database, got error: %s", err))
		return
	}

	r.updateModelFromAPI(ctx, &data, db)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

//nolint:dupl // database and pipeline resources have similar update patterns
func (r *DatabaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DatabaseResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := sevallaapi.UpdateDatabaseRequest{}

	if !data.Name.IsNull() {
		name := data.Name.ValueString()
		updateReq.Name = &name
	}

	if !data.Size.IsNull() {
		size := data.Size.ValueString()
		updateReq.Size = &size
	}

	if !data.Password.IsNull() {
		password := data.Password.ValueString()
		updateReq.Password = &password
	}

	db, err := sevallaapi.NewDatabaseService(r.client).Update(ctx, data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update database, got error: %s", err))
		return
	}

	r.updateModelFromAPI(ctx, &data, db)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DatabaseResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := sevallaapi.NewDatabaseService(r.client).Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete database, got error: %s", err))
		return
	}
}

func (r *DatabaseResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *DatabaseResource) updateModelFromAPI(
	_ context.Context,
	data *DatabaseResourceModel,
	db *sevallaapi.Database,
) {
	data.ID = types.StringValue(db.ID)
	data.Name = types.StringValue(db.Name)
	data.Type = types.StringValue(db.Type)
	data.Version = types.StringValue(db.Version)
	data.Size = types.StringValue(db.Size)
	data.Host = types.StringValue(db.Host)
	data.Port = types.Int64Value(int64(db.Port))
	data.Username = types.StringValue(db.Username)
	data.Status = types.StringValue(db.Status)
	data.CreatedAt = types.StringValue(db.CreatedAt.Format("2006-01-02T15:04:05Z"))
	data.UpdatedAt = types.StringValue(db.UpdatedAt.Format("2006-01-02T15:04:05Z"))

	// Only update password if it's not empty in the response
	if db.Password != "" {
		data.Password = types.StringValue(db.Password)
	}
}