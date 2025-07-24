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
var _ datasource.DataSource = &ApplicationDataSource{}

func NewApplicationDataSource() datasource.DataSource {
	return &ApplicationDataSource{}
}

// ApplicationDataSource defines the data source implementation.
type ApplicationDataSource struct {
	client *sevallaapi.Client
}

// ApplicationDataSourceModel describes the data source data model.
type ApplicationDataSourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	DisplayName          types.String `tfsdk:"display_name"`
	Status               types.String `tfsdk:"status"`
	CompanyID            types.String `tfsdk:"company_id"`
	RepoURL              types.String `tfsdk:"repo_url"`
	DefaultBranch        types.String `tfsdk:"default_branch"`
	AutoDeploy           types.Bool   `tfsdk:"auto_deploy"`
	BuildPath            types.String `tfsdk:"build_path"`
	BuildType            types.String `tfsdk:"build_type"`
	NodeVersion          types.String `tfsdk:"node_version"`
	DockerfilePath       types.String `tfsdk:"dockerfile_path"`
	DockerComposeFile    types.String `tfsdk:"docker_compose_file"`
	StartCommand         types.String `tfsdk:"start_command"`
	InstallCommand       types.String `tfsdk:"install_command"`
	EnvironmentVariables types.List   `tfsdk:"environment_variables"`
	CreatedAt            types.Int64  `tfsdk:"created_at"`
	UpdatedAt            types.Int64  `tfsdk:"updated_at"`
	Deployments          types.List   `tfsdk:"deployments"`
	Processes            types.List   `tfsdk:"processes"`
	InternalConnections  types.List   `tfsdk:"internal_connections"`
}

func (d *ApplicationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (d *ApplicationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches information about a specific Sevalla application including configuration details, deployments, processes, and internal connections.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier of the application.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique name of the application.",
			},
			"display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display name of the application.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The current status of the application (deploying, deployed, failed, stopped).",
			},
			"company_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The company ID that owns this application.",
			},
			"repo_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The repository URL for the application.",
			},
			"default_branch": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The default branch to deploy from.",
			},
			"auto_deploy": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether to automatically deploy on git push.",
			},
			"build_path": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The build path for the application.",
			},
			"build_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The build type (dockerfile, pack, nixpacks).",
			},
			"node_version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The Node.js version to use (16.20.0, 18.16.0, 20.2.0).",
			},
			"dockerfile_path": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The path to the Dockerfile.",
			},
			"docker_compose_file": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The path to the docker-compose file.",
			},
			"start_command": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The start command for the application.",
			},
			"install_command": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The install command for the application.",
			},
			"environment_variables": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Environment variables for the application.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The environment variable key.",
						},
						"value": schema.StringAttribute{
							Computed:            true,
							Sensitive:           true,
							MarkdownDescription: "The environment variable value.",
						},
					},
				},
			},
			"created_at": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The timestamp when the application was created.",
			},
			"updated_at": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The timestamp when the application was last updated.",
			},
			"deployments": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of deployments for this application.",
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
						"branch": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The branch for this deployment.",
						},
						"repo_url": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The repository URL.",
						},
						"commit_hash": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The commit hash.",
						},
						"commit_message": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The commit message.",
						},
						"created_at": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "When the deployment was created.",
						},
						"updated_at": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "When the deployment was last updated.",
						},
						"build_logs": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The build logs.",
						},
					},
				},
			},
			"processes": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of processes for this application.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The process ID.",
						},
						"key": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The process key.",
						},
						"type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The process type.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The process display name.",
						},
						"resource_type_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The resource type name.",
						},
						"entrypoint": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The process entrypoint.",
						},
					},
				},
			},
			"internal_connections": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of internal connections for this application.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The connection ID.",
						},
						"target_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The target type (appResource, dbResource, envResource).",
						},
						"target_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The target resource ID.",
						},
						"created_at": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "When the connection was created.",
						},
					},
				},
			},
		},
	}
}

func (d *ApplicationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ApplicationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ApplicationDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading application", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	app, err := d.client.Applications.Get(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read application, got error: %s", err))
		return
	}

	// Map all fields from API response using the same logic as the resource
	d.mapApplicationToModel(ctx, &data, &app.App)

	tflog.Trace(ctx, "Read application data source")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapApplicationToModel maps API response to Terraform model (same as resource)
func (d *ApplicationDataSource) mapApplicationToModel(ctx context.Context, data *ApplicationDataSourceModel, app *sevallaapi.ApplicationDetails) {
	data.ID = types.StringValue(app.ID)
	data.Name = types.StringValue(app.Name)
	data.DisplayName = types.StringValue(app.DisplayName)
	data.Status = types.StringValue(app.Status)
	data.CompanyID = types.StringValue(app.CompanyID)
	data.CreatedAt = types.Int64Value(app.CreatedAt)
	data.UpdatedAt = types.Int64Value(app.UpdatedAt)

	// Repository fields
	data.RepoURL = types.StringValue(app.RepoURL)
	data.DefaultBranch = types.StringValue(app.DefaultBranch)
	data.AutoDeploy = types.BoolValue(app.AutoDeploy)

	// Build configuration
	data.BuildPath = types.StringValue(app.BuildPath)
	data.BuildType = types.StringValue(app.BuildType)
	data.NodeVersion = types.StringValue(app.NodeVersion)
	data.DockerfilePath = types.StringValue(app.DockerfilePath)
	data.DockerComposeFile = types.StringValue(app.DockerComposeFile)
	data.StartCommand = types.StringValue(app.StartCommand)
	data.InstallCommand = types.StringValue(app.InstallCommand)

	// Convert environment variables
	envVars := make([]attr.Value, len(app.EnvironmentVariables))
	for i, envVar := range app.EnvironmentVariables {
		envVarObj, _ := types.ObjectValue(
			map[string]attr.Type{
				"key":   types.StringType,
				"value": types.StringType,
			},
			map[string]attr.Value{
				"key":   types.StringValue(envVar.Key),
				"value": types.StringValue(envVar.Value),
			},
		)
		envVars[i] = envVarObj
	}
	data.EnvironmentVariables, _ = types.ListValue(
		types.ObjectType{AttrTypes: map[string]attr.Type{"key": types.StringType, "value": types.StringType}},
		envVars,
	)

	// Convert deployments
	deployments := make([]attr.Value, len(app.Deployments))
	for i, deployment := range app.Deployments {
		commitMsg := ""
		if deployment.CommitMessage != nil {
			commitMsg = *deployment.CommitMessage
		}
		deploymentObj, _ := types.ObjectValue(
			map[string]attr.Type{
				"id":             types.StringType,
				"status":         types.StringType,
				"branch":         types.StringType,
				"repo_url":       types.StringType,
				"commit_hash":    types.StringType,
				"commit_message": types.StringType,
				"created_at":     types.Int64Type,
				"updated_at":     types.Int64Type,
				"build_logs":     types.StringType,
			},
			map[string]attr.Value{
				"id":             types.StringValue(deployment.ID),
				"status":         types.StringValue(deployment.Status),
				"branch":         types.StringValue(deployment.Branch),
				"repo_url":       types.StringValue(deployment.RepoURL),
				"commit_hash":    types.StringValue(deployment.CommitHash),
				"commit_message": types.StringValue(commitMsg),
				"created_at":     types.Int64Value(deployment.CreatedAt),
				"updated_at":     types.Int64Value(deployment.UpdatedAt),
				"build_logs":     types.StringValue(deployment.BuildLogs),
			},
		)
		deployments[i] = deploymentObj
	}
	deploymentAttrTypes := map[string]attr.Type{
		"id":             types.StringType,
		"status":         types.StringType,
		"branch":         types.StringType,
		"repo_url":       types.StringType,
		"commit_hash":    types.StringType,
		"commit_message": types.StringType,
		"created_at":     types.Int64Type,
		"updated_at":     types.Int64Type,
		"build_logs":     types.StringType,
	}
	data.Deployments, _ = types.ListValue(types.ObjectType{AttrTypes: deploymentAttrTypes}, deployments)

	// Convert processes
	processes := make([]attr.Value, len(app.Processes))
	for i, process := range app.Processes {
		processObj, _ := types.ObjectValue(
			map[string]attr.Type{
				"id":                 types.StringType,
				"key":                types.StringType,
				"type":               types.StringType,
				"display_name":       types.StringType,
				"resource_type_name": types.StringType,
				"entrypoint":         types.StringType,
			},
			map[string]attr.Value{
				"id":                 types.StringValue(process.ID),
				"key":                types.StringValue(process.Key),
				"type":               types.StringValue(process.Type),
				"display_name":       types.StringValue(process.DisplayName),
				"resource_type_name": types.StringValue(process.ResourceTypeName),
				"entrypoint":         types.StringValue(process.Entrypoint),
			},
		)
		processes[i] = processObj
	}
	processAttrTypes := map[string]attr.Type{
		"id":                 types.StringType,
		"key":                types.StringType,
		"type":               types.StringType,
		"display_name":       types.StringType,
		"resource_type_name": types.StringType,
		"entrypoint":         types.StringType,
	}
	data.Processes, _ = types.ListValue(types.ObjectType{AttrTypes: processAttrTypes}, processes)

	// Convert internal connections
	connections := make([]attr.Value, len(app.InternalConnections))
	for i, conn := range app.InternalConnections {
		connObj, _ := types.ObjectValue(
			map[string]attr.Type{
				"id":          types.StringType,
				"target_type": types.StringType,
				"target_id":   types.StringType,
				"created_at":  types.Int64Type,
			},
			map[string]attr.Value{
				"id":          types.StringValue(conn.ID),
				"target_type": types.StringValue(conn.TargetType),
				"target_id":   types.StringValue(conn.TargetID),
				"created_at":  types.Int64Value(conn.CreatedAt),
			},
		)
		connections[i] = connObj
	}
	connAttrTypes := map[string]attr.Type{
		"id":          types.StringType,
		"target_type": types.StringType,
		"target_id":   types.StringType,
		"created_at":  types.Int64Type,
	}
	data.InternalConnections, _ = types.ListValue(types.ObjectType{AttrTypes: connAttrTypes}, connections)
}
