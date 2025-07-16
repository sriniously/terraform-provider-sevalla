package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sriniously/terraform-provider-sevalla/internal/sevallaapi"
)

var _ resource.Resource = &ObjectStorageResource{}
var _ resource.ResourceWithImportState = &ObjectStorageResource{}

func NewObjectStorageResource() resource.Resource {
	return &ObjectStorageResource{}
}

type ObjectStorageResource struct {
	client *sevallaapi.Client
}

type ObjectStorageResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Region    types.String `tfsdk:"region"`
	Size      types.Int64  `tfsdk:"size"`
	Objects   types.Int64  `tfsdk:"objects"`
	Endpoint  types.String `tfsdk:"endpoint"`
	AccessKey types.String `tfsdk:"access_key"`
	SecretKey types.String `tfsdk:"secret_key"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

func (r *ObjectStorageResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_object_storage"
}

func (r *ObjectStorageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Sevalla object storage bucket.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Object storage identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Object storage bucket name",
				Required:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "Object storage region",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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

func (r *ObjectStorageResource) Configure(
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

func (r *ObjectStorageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ObjectStorageResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := sevallaapi.CreateObjectStorageRequest{
		Name: data.Name.ValueString(),
	}

	if !data.Region.IsNull() {
		createReq.Region = data.Region.ValueString()
	}

	tflog.Trace(ctx, "creating object storage")

	bucket, err := sevallaapi.NewObjectStorageService(r.client).Create(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create object storage, got error: %s", err))
		return
	}

	r.updateModelFromAPI(ctx, &data, bucket)

	tflog.Trace(ctx, "created object storage")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ObjectStorageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ObjectStorageResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucket, err := sevallaapi.NewObjectStorageService(r.client).Get(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read object storage, got error: %s", err))
		return
	}

	r.updateModelFromAPI(ctx, &data, bucket)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ObjectStorageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ObjectStorageResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := sevallaapi.UpdateObjectStorageRequest{}

	if !data.Name.IsNull() {
		name := data.Name.ValueString()
		updateReq.Name = &name
	}

	bucket, err := sevallaapi.NewObjectStorageService(r.client).Update(ctx, data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update object storage, got error: %s", err))
		return
	}

	r.updateModelFromAPI(ctx, &data, bucket)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ObjectStorageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ObjectStorageResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := sevallaapi.NewObjectStorageService(r.client).Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete object storage, got error: %s", err))
		return
	}
}

func (r *ObjectStorageResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ObjectStorageResource) updateModelFromAPI(
	_ context.Context,
	data *ObjectStorageResourceModel,
	bucket *sevallaapi.ObjectStorage,
) {
	data.ID = types.StringValue(bucket.ID)
	data.Name = types.StringValue(bucket.Name)
	data.Region = types.StringValue(bucket.Region)
	data.Size = types.Int64Value(bucket.Size)
	data.Objects = types.Int64Value(int64(bucket.Objects))
	data.Endpoint = types.StringValue(bucket.Endpoint)
	data.AccessKey = types.StringValue(bucket.AccessKey)
	data.SecretKey = types.StringValue(bucket.SecretKey)
	data.CreatedAt = types.StringValue(bucket.CreatedAt.Format("2006-01-02T15:04:05Z"))
	data.UpdatedAt = types.StringValue(bucket.UpdatedAt.Format("2006-01-02T15:04:05Z"))
}
