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

var _ resource.Resource = &PipelineResource{}
var _ resource.ResourceWithImportState = &PipelineResource{}

func NewPipelineResource() resource.Resource {
	return &PipelineResource{}
}

type PipelineResource struct {
	client *sevallaapi.Client
}

type PipelineResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	AppID      types.String `tfsdk:"app_id"`
	Branch     types.String `tfsdk:"branch"`
	AutoDeploy types.Bool   `tfsdk:"auto_deploy"`
	CreatedAt  types.String `tfsdk:"created_at"`
	UpdatedAt  types.String `tfsdk:"updated_at"`
}

func (r *PipelineResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_pipeline"
}

func (r *PipelineResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Sevalla deployment pipeline.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Pipeline identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Pipeline name",
				Required:            true,
			},
			"app_id": schema.StringAttribute{
				MarkdownDescription: "Associated application ID",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"branch": schema.StringAttribute{
				MarkdownDescription: "Git branch for the pipeline",
				Required:            true,
			},
			"auto_deploy": schema.BoolAttribute{
				MarkdownDescription: "Enable automatic deployment",
				Optional:            true,
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

func (r *PipelineResource) Configure(
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

func (r *PipelineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PipelineResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := sevallaapi.CreatePipelineRequest{
		Name:   data.Name.ValueString(),
		AppID:  data.AppID.ValueString(),
		Branch: data.Branch.ValueString(),
	}

	if !data.AutoDeploy.IsNull() {
		createReq.AutoDeploy = data.AutoDeploy.ValueBool()
	}

	tflog.Trace(ctx, "creating pipeline")

	pipeline, err := sevallaapi.NewPipelineService(r.client).Create(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create pipeline, got error: %s", err))
		return
	}

	r.updateModelFromAPI(ctx, &data, pipeline)

	tflog.Trace(ctx, "created pipeline")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PipelineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PipelineResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pipeline, err := sevallaapi.NewPipelineService(r.client).Get(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read pipeline, got error: %s", err))
		return
	}

	r.updateModelFromAPI(ctx, &data, pipeline)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

//nolint:dupl // database and pipeline resources have similar update patterns
func (r *PipelineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PipelineResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := sevallaapi.UpdatePipelineRequest{}

	if !data.Name.IsNull() {
		name := data.Name.ValueString()
		updateReq.Name = &name
	}

	if !data.Branch.IsNull() {
		branch := data.Branch.ValueString()
		updateReq.Branch = &branch
	}

	if !data.AutoDeploy.IsNull() {
		autoDeploy := data.AutoDeploy.ValueBool()
		updateReq.AutoDeploy = &autoDeploy
	}

	pipeline, err := sevallaapi.NewPipelineService(r.client).Update(ctx, data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update pipeline, got error: %s", err))
		return
	}

	r.updateModelFromAPI(ctx, &data, pipeline)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PipelineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PipelineResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := sevallaapi.NewPipelineService(r.client).Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete pipeline, got error: %s", err))
		return
	}
}

func (r *PipelineResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *PipelineResource) updateModelFromAPI(
	_ context.Context,
	data *PipelineResourceModel,
	pipeline *sevallaapi.Pipeline,
) {
	data.ID = types.StringValue(pipeline.ID)
	data.Name = types.StringValue(pipeline.Name)
	data.AppID = types.StringValue(pipeline.AppID)
	data.Branch = types.StringValue(pipeline.Branch)
	data.AutoDeploy = types.BoolValue(pipeline.AutoDeploy)
	data.CreatedAt = types.StringValue(pipeline.CreatedAt.Format("2006-01-02T15:04:05Z"))
	data.UpdatedAt = types.StringValue(pipeline.UpdatedAt.Format("2006-01-02T15:04:05Z"))
}