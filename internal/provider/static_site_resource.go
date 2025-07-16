package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sriniously/terraform-provider-sevalla/internal/sevallaapi"
)

var _ resource.Resource = &StaticSiteResource{}
var _ resource.ResourceWithImportState = &StaticSiteResource{}

func NewStaticSiteResource() resource.Resource {
	return &StaticSiteResource{}
}

type StaticSiteResource struct {
	client *sevallaapi.Client
}

type StaticSiteResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Domain     types.String `tfsdk:"domain"`
	Repository types.Object `tfsdk:"repository"`
	Branch     types.String `tfsdk:"branch"`
	BuildDir   types.String `tfsdk:"build_dir"`
	BuildCmd   types.String `tfsdk:"build_cmd"`
	Status     types.String `tfsdk:"status"`
	CreatedAt  types.String `tfsdk:"created_at"`
	UpdatedAt  types.String `tfsdk:"updated_at"`
}

func (r *StaticSiteResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_static_site"
}

func (r *StaticSiteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Sevalla static site.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Static site identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Static site name",
				Required:            true,
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "Custom domain for the static site",
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
			"build_dir": schema.StringAttribute{
				MarkdownDescription: "Build output directory",
				Optional:            true,
			},
			"build_cmd": schema.StringAttribute{
				MarkdownDescription: "Build command to run",
				Optional:            true,
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

func (r *StaticSiteResource) Configure(
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

func (r *StaticSiteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data StaticSiteResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := sevallaapi.CreateStaticSiteRequest{
		Name: data.Name.ValueString(),
	}

	if !data.Branch.IsNull() {
		createReq.Branch = data.Branch.ValueString()
	}

	if !data.BuildDir.IsNull() {
		createReq.BuildDir = data.BuildDir.ValueString()
	}

	if !data.BuildCmd.IsNull() {
		createReq.BuildCmd = data.BuildCmd.ValueString()
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

	tflog.Trace(ctx, "creating static site")

	site, err := sevallaapi.NewStaticSiteService(r.client).Create(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create static site, got error: %s", err))
		return
	}

	r.updateModelFromAPI(ctx, &data, site)

	tflog.Trace(ctx, "created static site")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StaticSiteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data StaticSiteResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	site, err := sevallaapi.NewStaticSiteService(r.client).Get(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read static site, got error: %s", err))
		return
	}

	r.updateModelFromAPI(ctx, &data, site)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StaticSiteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data StaticSiteResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := sevallaapi.UpdateStaticSiteRequest{}

	if !data.Name.IsNull() {
		name := data.Name.ValueString()
		updateReq.Name = &name
	}

	if !data.Branch.IsNull() {
		branch := data.Branch.ValueString()
		updateReq.Branch = &branch
	}

	if !data.BuildDir.IsNull() {
		buildDir := data.BuildDir.ValueString()
		updateReq.BuildDir = &buildDir
	}

	if !data.BuildCmd.IsNull() {
		buildCmd := data.BuildCmd.ValueString()
		updateReq.BuildCmd = &buildCmd
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

	site, err := sevallaapi.NewStaticSiteService(r.client).Update(ctx, data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update static site, got error: %s", err))
		return
	}

	r.updateModelFromAPI(ctx, &data, site)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StaticSiteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data StaticSiteResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := sevallaapi.NewStaticSiteService(r.client).Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete static site, got error: %s", err))
		return
	}
}

func (r *StaticSiteResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *StaticSiteResource) updateModelFromAPI(
	_ context.Context,
	data *StaticSiteResourceModel,
	site *sevallaapi.StaticSite,
) {
	data.ID = types.StringValue(site.ID)
	data.Name = types.StringValue(site.Name)
	data.Domain = types.StringValue(site.Domain)
	data.Branch = types.StringValue(site.Branch)
	data.BuildDir = types.StringValue(site.BuildDir)
	data.BuildCmd = types.StringValue(site.BuildCmd)
	data.Status = types.StringValue(site.Status)
	data.CreatedAt = types.StringValue(site.CreatedAt.Format("2006-01-02T15:04:05Z"))
	data.UpdatedAt = types.StringValue(site.UpdatedAt.Format("2006-01-02T15:04:05Z"))

	if site.Repository != nil {
		repoObj := map[string]attr.Value{
			"url":    types.StringValue(site.Repository.URL),
			"type":   types.StringValue(site.Repository.Type),
			"branch": types.StringValue(site.Repository.Branch),
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