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
var _ resource.Resource = &StaticSiteResource{}
var _ resource.ResourceWithImportState = &StaticSiteResource{}

func NewStaticSiteResource() resource.Resource {
	return &StaticSiteResource{}
}

// StaticSiteResource defines the resource implementation.
type StaticSiteResource struct {
	client *sevallaapi.Client
}

// StaticSiteResourceModel describes the resource data model.
type StaticSiteResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	DisplayName        types.String `tfsdk:"display_name"`
	CompanyID          types.String `tfsdk:"company_id"`
	Status             types.String `tfsdk:"status"`
	RepoURL            types.String `tfsdk:"repo_url"`
	DefaultBranch      types.String `tfsdk:"default_branch"`
	AutoDeploy         types.Bool   `tfsdk:"auto_deploy"`
	GitType            types.String `tfsdk:"git_type"`
	Hostname           types.String `tfsdk:"hostname"`
	BuildCommand       types.String `tfsdk:"build_command"`
	NodeVersion        types.String `tfsdk:"node_version"`
	PublishedDirectory types.String `tfsdk:"published_directory"`
}

func (r *StaticSiteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_static_site"
}

func (r *StaticSiteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a static site on Sevalla platform.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the static site.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique name of the static site.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name of the static site.",
			},
			"company_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The company ID that owns this static site.",
			},
			"repo_url": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The repository URL for the static site.",
			},
			"default_branch": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The default branch to deploy from.",
			},
			"auto_deploy": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to automatically deploy on git push.",
			},
			"build_command": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The build command to run.",
			},
			"node_version": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The Node.js version to use (16.20.0, 18.16.0, 20.2.0).",
				Validators: []validator.String{
					stringvalidator.OneOf("16.20.0", "18.16.0", "20.2.0"),
				},
			},
			"published_directory": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The directory containing the built static files.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The current status of the static site.",
			},
			"git_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The git provider type (github, gitlab, etc.).",
			},
			"hostname": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The hostname where the static site is deployed.",
			},
		},
	}
}

func (r *StaticSiteResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *StaticSiteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data StaticSiteResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := sevallaapi.CreateStaticSiteRequest{
		CompanyID:   data.CompanyID.ValueString(),
		DisplayName: data.DisplayName.ValueString(),
		RepoURL:     data.RepoURL.ValueString(),
	}

	if !data.DefaultBranch.IsNull() {
		branch := data.DefaultBranch.ValueString()
		createReq.Branch = &branch
	}

	tflog.Debug(ctx, "Creating static site", map[string]interface{}{
		"company_id":   createReq.CompanyID,
		"display_name": createReq.DisplayName,
		"repo_url":     createReq.RepoURL,
	})

	site, err := r.client.StaticSites.Create(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create static site, got error: %s", err))
		return
	}

	data.ID = types.StringValue(site.StaticSite.ID)
	data.Name = types.StringValue(site.StaticSite.Name)
	data.DisplayName = types.StringValue(site.StaticSite.DisplayName)
	data.Status = types.StringValue(site.StaticSite.Status)
	data.RepoURL = types.StringValue(site.StaticSite.RepoURL)
	data.DefaultBranch = types.StringValue(site.StaticSite.DefaultBranch)
	data.AutoDeploy = types.BoolValue(site.StaticSite.AutoDeploy)
	data.GitType = types.StringValue(site.StaticSite.GitType)
	data.Hostname = types.StringValue(site.StaticSite.Hostname)

	if site.StaticSite.BuildCommand != nil {
		data.BuildCommand = types.StringValue(*site.StaticSite.BuildCommand)
	}

	tflog.Trace(ctx, "Created static site resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StaticSiteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data StaticSiteResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	site, err := r.client.StaticSites.Get(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read static site, got error: %s", err))
		return
	}

	data.ID = types.StringValue(site.StaticSite.ID)
	data.Name = types.StringValue(site.StaticSite.Name)
	data.DisplayName = types.StringValue(site.StaticSite.DisplayName)
	data.Status = types.StringValue(site.StaticSite.Status)
	data.RepoURL = types.StringValue(site.StaticSite.RepoURL)
	data.DefaultBranch = types.StringValue(site.StaticSite.DefaultBranch)
	data.AutoDeploy = types.BoolValue(site.StaticSite.AutoDeploy)
	data.GitType = types.StringValue(site.StaticSite.GitType)
	data.Hostname = types.StringValue(site.StaticSite.Hostname)

	if site.StaticSite.BuildCommand != nil {
		data.BuildCommand = types.StringValue(*site.StaticSite.BuildCommand)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StaticSiteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data StaticSiteResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := sevallaapi.UpdateStaticSiteRequest{}

	if !data.DisplayName.IsNull() {
		updateReq.DisplayName = stringPointer(data.DisplayName.ValueString())
	}

	if !data.AutoDeploy.IsNull() {
		autoDeploy := data.AutoDeploy.ValueBool()
		updateReq.AutoDeploy = &autoDeploy
	}

	if !data.DefaultBranch.IsNull() {
		updateReq.DefaultBranch = stringPointer(data.DefaultBranch.ValueString())
	}

	if !data.BuildCommand.IsNull() {
		updateReq.BuildCommand = stringPointer(data.BuildCommand.ValueString())
	}

	if !data.NodeVersion.IsNull() {
		updateReq.NodeVersion = stringPointer(data.NodeVersion.ValueString())
	}

	if !data.PublishedDirectory.IsNull() {
		updateReq.PublishedDirectory = stringPointer(data.PublishedDirectory.ValueString())
	}

	site, err := r.client.StaticSites.Update(ctx, data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update static site, got error: %s", err))
		return
	}

	data.ID = types.StringValue(site.StaticSite.ID)
	data.Name = types.StringValue(site.StaticSite.Name)
	data.DisplayName = types.StringValue(site.StaticSite.DisplayName)
	data.Status = types.StringValue(site.StaticSite.Status)
	data.AutoDeploy = types.BoolValue(site.StaticSite.AutoDeploy)
	data.DefaultBranch = types.StringValue(site.StaticSite.DefaultBranch)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StaticSiteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data StaticSiteResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.StaticSites.Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete static site, got error: %s", err))
		return
	}
}

func (r *StaticSiteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
