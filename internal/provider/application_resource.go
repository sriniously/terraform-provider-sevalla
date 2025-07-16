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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sriniously/terraform-provider-sevalla/internal/sevallaapi"
)

var _ resource.Resource = &ApplicationResource{}
var _ resource.ResourceWithImportState = &ApplicationResource{}

func NewApplicationResource() resource.Resource {
	return &ApplicationResource{}
}

type ApplicationResource struct {
	client *sevallaapi.Client
}

type ApplicationResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Domain       types.String `tfsdk:"domain"`
	Repository   types.Object `tfsdk:"repository"`
	Branch       types.String `tfsdk:"branch"`
	BuildCommand types.String `tfsdk:"build_command"`
	StartCommand types.String `tfsdk:"start_command"`
	Environment  types.Map    `tfsdk:"environment"`
	Instances    types.Int64  `tfsdk:"instances"`
	Memory       types.Int64  `tfsdk:"memory"`
	CPU          types.Int64  `tfsdk:"cpu"`
	Status       types.String `tfsdk:"status"`
	CreatedAt    types.String `tfsdk:"created_at"`
	UpdatedAt    types.String `tfsdk:"updated_at"`
}

type RepositoryModel struct {
	URL    types.String `tfsdk:"url"`
	Type   types.String `tfsdk:"type"`
	Branch types.String `tfsdk:"branch"`
}

func (r *ApplicationResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (r *ApplicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Sevalla application.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Application identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Application name",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Application description",
				Optional:            true,
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "Custom domain for the application",
				Optional:            true,
			},
			"repository": schema.SingleNestedAttribute{
				MarkdownDescription: "Source code repository configuration",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"url": schema.StringAttribute{
						MarkdownDescription: "Repository URL",
						Required:            true,
					},
					"type": schema.StringAttribute{
						MarkdownDescription: "Repository type (github, gitlab, bitbucket)",
						Required:            true,
					},
					"branch": schema.StringAttribute{
						MarkdownDescription: "Repository branch",
						Optional:            true,
					},
				},
			},
			"branch": schema.StringAttribute{
				MarkdownDescription: "Git branch to deploy",
				Optional:            true,
			},
			"build_command": schema.StringAttribute{
				MarkdownDescription: "Build command to run",
				Optional:            true,
			},
			"start_command": schema.StringAttribute{
				MarkdownDescription: "Start command to run",
				Optional:            true,
			},
			"environment": schema.MapAttribute{
				MarkdownDescription: "Environment variables",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"instances": schema.Int64Attribute{
				MarkdownDescription: "Number of instances",
				Optional:            true,
			},
			"memory": schema.Int64Attribute{
				MarkdownDescription: "Memory allocation in MB",
				Optional:            true,
			},
			"cpu": schema.Int64Attribute{
				MarkdownDescription: "CPU allocation in millicores",
				Optional:            true,
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

func (r *ApplicationResource) Configure(
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

//nolint:cyclop // terraform resource methods require handling multiple conditional fields
func (r *ApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ApplicationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert Terraform model to API request
	createReq := sevallaapi.CreateApplicationRequest{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
	}

	if !data.Branch.IsNull() {
		createReq.Branch = data.Branch.ValueString()
	}

	if !data.BuildCommand.IsNull() {
		createReq.BuildCommand = data.BuildCommand.ValueString()
	}

	if !data.StartCommand.IsNull() {
		createReq.StartCommand = data.StartCommand.ValueString()
	}

	if !data.Environment.IsNull() {
		env := make(map[string]string)
		resp.Diagnostics.Append(data.Environment.ElementsAs(ctx, &env, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.Environment = env
	}

	if !data.Instances.IsNull() {
		createReq.Instances = int(data.Instances.ValueInt64())
	}

	if !data.Memory.IsNull() {
		createReq.Memory = int(data.Memory.ValueInt64())
	}

	if !data.CPU.IsNull() {
		createReq.CPU = int(data.CPU.ValueInt64())
	}

	// Handle repository
	if !data.Repository.IsNull() {
		var repo RepositoryModel
		resp.Diagnostics.Append(data.Repository.As(ctx, &repo, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}

		createReq.Repository = &sevallaapi.Repository{
			URL:    repo.URL.ValueString(),
			Type:   repo.Type.ValueString(),
			Branch: repo.Branch.ValueString(),
		}
	}

	tflog.Trace(ctx, "creating application")

	app, err := sevallaapi.NewApplicationService(r.client).Create(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create application, got error: %s", err))
		return
	}

	// Update the state with the created application
	r.updateModelFromAPI(ctx, &data, app)

	tflog.Trace(ctx, "created application")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ApplicationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	app, err := sevallaapi.NewApplicationService(r.client).Get(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read application, got error: %s", err))
		return
	}

	r.updateModelFromAPI(ctx, &data, app)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

//nolint:cyclop // terraform resource methods require handling multiple conditional fields
func (r *ApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ApplicationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert Terraform model to API request
	updateReq := sevallaapi.UpdateApplicationRequest{}

	if !data.Name.IsNull() {
		name := data.Name.ValueString()
		updateReq.Name = &name
	}

	if !data.Description.IsNull() {
		desc := data.Description.ValueString()
		updateReq.Description = &desc
	}

	if !data.Branch.IsNull() {
		branch := data.Branch.ValueString()
		updateReq.Branch = &branch
	}

	if !data.BuildCommand.IsNull() {
		buildCmd := data.BuildCommand.ValueString()
		updateReq.BuildCommand = &buildCmd
	}

	if !data.StartCommand.IsNull() {
		startCmd := data.StartCommand.ValueString()
		updateReq.StartCommand = &startCmd
	}

	if !data.Environment.IsNull() {
		env := make(map[string]string)
		resp.Diagnostics.Append(data.Environment.ElementsAs(ctx, &env, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateReq.Environment = &env
	}

	if !data.Instances.IsNull() {
		instances := int(data.Instances.ValueInt64())
		updateReq.Instances = &instances
	}

	if !data.Memory.IsNull() {
		memory := int(data.Memory.ValueInt64())
		updateReq.Memory = &memory
	}

	if !data.CPU.IsNull() {
		cpu := int(data.CPU.ValueInt64())
		updateReq.CPU = &cpu
	}

	// Handle repository
	if !data.Repository.IsNull() {
		var repo RepositoryModel
		resp.Diagnostics.Append(data.Repository.As(ctx, &repo, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}

		updateReq.Repository = &sevallaapi.Repository{
			URL:    repo.URL.ValueString(),
			Type:   repo.Type.ValueString(),
			Branch: repo.Branch.ValueString(),
		}
	}

	app, err := sevallaapi.NewApplicationService(r.client).Update(ctx, data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update application, got error: %s", err))
		return
	}

	r.updateModelFromAPI(ctx, &data, app)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ApplicationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := sevallaapi.NewApplicationService(r.client).Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete application, got error: %s", err))
		return
	}
}

func (r *ApplicationResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ApplicationResource) updateModelFromAPI(
	_ context.Context,
	data *ApplicationResourceModel,
	app *sevallaapi.Application,
) {
	data.ID = types.StringValue(app.ID)
	data.Name = types.StringValue(app.Name)
	data.Description = types.StringValue(app.Description)
	data.Domain = types.StringValue(app.Domain)
	data.Branch = types.StringValue(app.Branch)
	data.BuildCommand = types.StringValue(app.BuildCommand)
	data.StartCommand = types.StringValue(app.StartCommand)
	data.Status = types.StringValue(app.Status)
	data.CreatedAt = types.StringValue(app.CreatedAt.Format("2006-01-02T15:04:05Z"))
	data.UpdatedAt = types.StringValue(app.UpdatedAt.Format("2006-01-02T15:04:05Z"))

	if app.Environment != nil {
		envMap := make(map[string]attr.Value)
		for k, v := range app.Environment {
			envMap[k] = types.StringValue(v)
		}
		envValue, _ := types.MapValue(types.StringType, envMap)
		data.Environment = envValue
	}

	if app.Instances > 0 {
		data.Instances = types.Int64Value(int64(app.Instances))
	}

	if app.Memory > 0 {
		data.Memory = types.Int64Value(int64(app.Memory))
	}

	if app.CPU > 0 {
		data.CPU = types.Int64Value(int64(app.CPU))
	}

	if app.Repository != nil {
		repoObj := map[string]attr.Value{
			"url":    types.StringValue(app.Repository.URL),
			"type":   types.StringValue(app.Repository.Type),
			"branch": types.StringValue(app.Repository.Branch),
		}
		objValue, _ := types.ObjectValue(
			map[string]attr.Type{
				"url":    types.StringType,
				"type":   types.StringType,
				"branch": types.StringType,
			},
			repoObj,
		)
		data.Repository = objValue
	}
}