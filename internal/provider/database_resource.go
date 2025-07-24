package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sriniously/terraform-provider-sevalla/internal/sevallaapi"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DatabaseResource{}
var _ resource.ResourceWithImportState = &DatabaseResource{}

func NewDatabaseResource() resource.Resource {
	return &DatabaseResource{}
}

// DatabaseResource defines the resource implementation.
type DatabaseResource struct {
	client *sevallaapi.Client
}

// DatabaseResourceModel describes the resource data model.
type DatabaseResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	DisplayName      types.String `tfsdk:"display_name"`
	CompanyID        types.String `tfsdk:"company_id"`
	Location         types.String `tfsdk:"location"`
	ResourceType     types.String `tfsdk:"resource_type"`
	Type             types.String `tfsdk:"type"`
	Version          types.String `tfsdk:"version"`
	DBName           types.String `tfsdk:"db_name"`
	DBPassword       types.String `tfsdk:"db_password"`
	DBUser           types.String `tfsdk:"db_user"`
	Status           types.String `tfsdk:"status"`
	InternalHostname types.String `tfsdk:"internal_hostname"`
	InternalPort     types.String `tfsdk:"internal_port"`
	ExternalHostname types.String `tfsdk:"external_hostname"`
	ExternalPort     types.String `tfsdk:"external_port"`
}

func (r *DatabaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database"
}

func (r *DatabaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a database on Sevalla platform.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the database.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique name of the database.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name of the database.",
			},
			"company_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The company ID that owns this database.",
			},
			"location": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The location where the database will be created (e.g., us-central1, europe-west3).",
			},
			"resource_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The resource type for the database (db1, db2, ..., db9).",
				Validators: []validator.String{
					stringvalidator.OneOf("db1", "db2", "db3", "db4", "db5", "db6", "db7", "db8", "db9"),
				},
			},
			"type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The database type (postgresql, redis, mariadb, mysql).",
				Validators: []validator.String{
					stringvalidator.OneOf("postgresql", "redis", "mariadb", "mysql"),
				},
			},
			"version": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The database version.",
			},
			"db_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The database name.",
			},
			"db_password": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "The database password.",
			},
			"db_user": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The database user (optional for Redis, required for others).",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The current status of the database.",
			},
			"internal_hostname": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The internal hostname for database connections.",
			},
			"internal_port": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The internal port for database connections.",
			},
			"external_hostname": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The external hostname for database connections.",
			},
			"external_port": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The external port for database connections.",
			},
		},
	}
}

func (r *DatabaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(SevallaProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected SevallaProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = data.Client
}

func (r *DatabaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DatabaseResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := sevallaapi.CreateDatabaseRequest{
		CompanyID:    data.CompanyID.ValueString(),
		Location:     data.Location.ValueString(),
		ResourceType: data.ResourceType.ValueString(),
		DisplayName:  data.DisplayName.ValueString(),
		DBName:       data.DBName.ValueString(),
		DBPassword:   data.DBPassword.ValueString(),
		Type:         data.Type.ValueString(),
		Version:      data.Version.ValueString(),
	}

	if !data.DBUser.IsNull() {
		createReq.DBUser = data.DBUser.ValueString()
	}

	tflog.Debug(ctx, "Creating database", map[string]interface{}{
		"company_id":    createReq.CompanyID,
		"display_name":  createReq.DisplayName,
		"type":          createReq.Type,
		"version":       createReq.Version,
		"location":      createReq.Location,
		"resource_type": createReq.ResourceType,
	})

	db, err := r.client.Databases.Create(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create database, got error: %s", err))
		return
	}

	data.ID = types.StringValue(db.Database.ID)
	data.Name = types.StringValue(db.Database.Name)
	data.DisplayName = types.StringValue(db.Database.DisplayName)
	data.Status = types.StringValue(db.Database.Status)
	data.Type = types.StringValue(db.Database.Type)
	data.Version = types.StringValue(db.Database.Version)

	if db.Database.InternalHostname != nil {
		data.InternalHostname = types.StringValue(*db.Database.InternalHostname)
	} else {
		data.InternalHostname = types.StringNull()
	}
	if db.Database.InternalPort != nil {
		data.InternalPort = types.StringValue(*db.Database.InternalPort)
	} else {
		data.InternalPort = types.StringNull()
	}
	if db.Database.ExternalHostname != nil {
		data.ExternalHostname = types.StringValue(*db.Database.ExternalHostname)
	} else {
		data.ExternalHostname = types.StringNull()
	}
	if db.Database.ExternalPort != nil {
		data.ExternalPort = types.StringValue(*db.Database.ExternalPort)
	} else {
		data.ExternalPort = types.StringNull()
	}

	tflog.Trace(ctx, "Created database resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DatabaseResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	db, err := r.client.Databases.Get(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read database, got error: %s", err))
		return
	}

	data.ID = types.StringValue(db.Database.ID)
	data.Name = types.StringValue(db.Database.Name)
	data.DisplayName = types.StringValue(db.Database.DisplayName)
	data.Status = types.StringValue(db.Database.Status)
	data.Type = types.StringValue(db.Database.Type)
	data.Version = types.StringValue(db.Database.Version)

	if db.Database.InternalHostname != nil {
		data.InternalHostname = types.StringValue(*db.Database.InternalHostname)
	} else {
		data.InternalHostname = types.StringNull()
	}
	if db.Database.InternalPort != nil {
		data.InternalPort = types.StringValue(*db.Database.InternalPort)
	} else {
		data.InternalPort = types.StringNull()
	}
	if db.Database.ExternalHostname != nil {
		data.ExternalHostname = types.StringValue(*db.Database.ExternalHostname)
	} else {
		data.ExternalHostname = types.StringNull()
	}
	if db.Database.ExternalPort != nil {
		data.ExternalPort = types.StringValue(*db.Database.ExternalPort)
	} else {
		data.ExternalPort = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DatabaseResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := sevallaapi.UpdateDatabaseRequest{
		DisplayName: stringPointer(data.DisplayName.ValueString()),
	}

	if !data.ResourceType.IsNull() {
		updateReq.ResourceType = stringPointer(data.ResourceType.ValueString())
	}

	db, err := r.client.Databases.Update(ctx, data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update database, got error: %s", err))
		return
	}

	data.ID = types.StringValue(db.Database.ID)
	data.Name = types.StringValue(db.Database.Name)
	data.DisplayName = types.StringValue(db.Database.DisplayName)
	data.Status = types.StringValue(db.Database.Status)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DatabaseResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Databases.Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete database, got error: %s", err))
		return
	}
}

func (r *DatabaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
