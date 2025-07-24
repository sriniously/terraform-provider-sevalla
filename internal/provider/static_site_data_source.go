package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sriniously/terraform-provider-sevalla/internal/sevallaapi"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &StaticSiteDataSource{}

func NewStaticSiteDataSource() datasource.DataSource {
	return &StaticSiteDataSource{}
}

// StaticSiteDataSource defines the data source implementation.
type StaticSiteDataSource struct {
	client *sevallaapi.Client
}

// StaticSiteDeploymentModel represents a static site deployment.
type StaticSiteDeploymentModel struct {
	ID            types.String `tfsdk:"id"`
	Status        types.String `tfsdk:"status"`
	RepoURL       types.String `tfsdk:"repo_url"`
	Branch        types.String `tfsdk:"branch"`
	CommitMessage types.String `tfsdk:"commit_message"`
	CreatedAt     types.Int64  `tfsdk:"created_at"`
}

// StaticSiteDataSourceModel describes the data source data model.
type StaticSiteDataSourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	DisplayName        types.String `tfsdk:"display_name"`
	CompanyID          types.String `tfsdk:"company_id"`
	Status             types.String `tfsdk:"status"`
	RepoURL            types.String `tfsdk:"repo_url"`
	DefaultBranch      types.String `tfsdk:"default_branch"`
	AutoDeploy         types.Bool   `tfsdk:"auto_deploy"`
	RemoteRepositoryID types.String `tfsdk:"remote_repository_id"`
	GitRepositoryID    types.String `tfsdk:"git_repository_id"`
	GitType            types.String `tfsdk:"git_type"`
	Hostname           types.String `tfsdk:"hostname"`
	BuildCommand       types.String `tfsdk:"build_command"`
	CreatedAt          types.Int64  `tfsdk:"created_at"`
	UpdatedAt          types.Int64  `tfsdk:"updated_at"`
	Deployments        types.List   `tfsdk:"deployments"`
}

func (d *StaticSiteDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_static_site"
}

func (d *StaticSiteDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches information about a specific Sevalla static site including deployment history and configuration details.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier of the static site.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique name of the static site.",
			},
			"display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display name of the static site.",
			},
			"company_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The company ID that owns this static site.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The current status of the static site.",
			},
			"repo_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The repository URL for the static site.",
			},
			"default_branch": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The default branch to deploy from.",
			},
			"auto_deploy": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether to automatically deploy on git push.",
			},
			"remote_repository_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The remote repository identifier.",
			},
			"git_repository_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The git repository identifier.",
			},
			"git_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The git provider type (github, gitlab, etc.).",
			},
			"hostname": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The hostname where the static site is deployed.",
			},
			"build_command": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The build command used for the static site.",
			},
			"created_at": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The timestamp when the static site was created.",
			},
			"updated_at": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The timestamp when the static site was last updated.",
			},
			"deployments": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of deployments for this static site.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The deployment ID.",
						},
						"status": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The deployment status.",
						},
						"repo_url": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The repository URL.",
						},
						"branch": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The branch for this deployment.",
						},
						"commit_message": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The commit message.",
						},
						"created_at": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "When the deployment was created.",
						},
					},
				},
			},
		},
	}
}

func (d *StaticSiteDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(SevallaProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected SevallaProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = data.Client
}

func (d *StaticSiteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data StaticSiteDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading static site", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	site, err := d.client.StaticSites.Get(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read static site, got error: %s", err))
		return
	}

	// Map all fields from API response
	data.ID = types.StringValue(site.StaticSite.ID)
	data.Name = types.StringValue(site.StaticSite.Name)
	data.DisplayName = types.StringValue(site.StaticSite.DisplayName)
	data.Status = types.StringValue(site.StaticSite.Status)
	data.RepoURL = types.StringValue(site.StaticSite.RepoURL)
	data.DefaultBranch = types.StringValue(site.StaticSite.DefaultBranch)
	data.AutoDeploy = types.BoolValue(site.StaticSite.AutoDeploy)
	data.RemoteRepositoryID = types.StringValue(site.StaticSite.RemoteRepositoryID)
	data.GitRepositoryID = types.StringValue(site.StaticSite.GitRepositoryID)
	data.GitType = types.StringValue(site.StaticSite.GitType)
	data.Hostname = types.StringValue(site.StaticSite.Hostname)
	data.CreatedAt = types.Int64Value(site.StaticSite.CreatedAt)
	data.UpdatedAt = types.Int64Value(site.StaticSite.UpdatedAt)

	if site.StaticSite.BuildCommand != nil {
		data.BuildCommand = types.StringValue(*site.StaticSite.BuildCommand)
	} else {
		data.BuildCommand = types.StringNull()
	}

	// Convert deployments
	deployments := make([]attr.Value, len(site.StaticSite.Deployments))
	for i, deployment := range site.StaticSite.Deployments {
		commitMsg := ""
		if deployment.CommitMessage != nil {
			commitMsg = *deployment.CommitMessage
		}
		deploymentObj, _ := types.ObjectValue(
			map[string]attr.Type{
				"id":             types.StringType,
				"status":         types.StringType,
				"repo_url":       types.StringType,
				"branch":         types.StringType,
				"commit_message": types.StringType,
				"created_at":     types.Int64Type,
			},
			map[string]attr.Value{
				"id":             types.StringValue(deployment.ID),
				"status":         types.StringValue(deployment.Status),
				"repo_url":       types.StringValue(deployment.RepoURL),
				"branch":         types.StringValue(deployment.Branch),
				"commit_message": types.StringValue(commitMsg),
				"created_at":     types.Int64Value(deployment.CreatedAt),
			},
		)
		deployments[i] = deploymentObj
	}
	deploymentAttrTypes := map[string]attr.Type{
		"id":             types.StringType,
		"status":         types.StringType,
		"repo_url":       types.StringType,
		"branch":         types.StringType,
		"commit_message": types.StringType,
		"created_at":     types.Int64Type,
	}
	data.Deployments, _ = types.ListValue(types.ObjectType{AttrTypes: deploymentAttrTypes}, deployments)

	tflog.Trace(ctx, "Read static site data source")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
