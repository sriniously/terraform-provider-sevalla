package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sriniously/terraform-provider-sevalla/internal/sevallaapi"
)

var _ provider.Provider = &SevallaProvider{}

type SevallaProvider struct {
	version string
}

type SevallaProviderModel struct {
	Token   types.String `tfsdk:"token"`
	BaseURL types.String `tfsdk:"base_url"`
}

type SevallaProviderData struct {
	Client *sevallaapi.Client
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SevallaProvider{
			version: version,
		}
	}
}

func (p *SevallaProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sevalla"
	resp.Version = p.version
}

func (p *SevallaProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				MarkdownDescription: "The Sevalla API token. Can also be set via the `SEVALLA_TOKEN` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "The base URL for the Sevalla API. Defaults to `https://api.sevalla.com`.",
				Optional:            true,
			},
		},
	}
}

func (p *SevallaProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse,
) {
	var data SevallaProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values
	token := os.Getenv("SEVALLA_TOKEN")
	baseURL := sevallaapi.DefaultBaseURL

	if !data.Token.IsNull() {
		token = data.Token.ValueString()
	}

	if !data.BaseURL.IsNull() {
		baseURL = data.BaseURL.ValueString()
	}

	// Check if token is provided
	if token == "" {
		resp.Diagnostics.AddError(
			"Unable to find token",
			"Token cannot be an empty string. Please set the token in the provider "+
				"configuration or via the SEVALLA_TOKEN environment variable.",
		)
		return
	}

	ctx = tflog.SetField(ctx, "sevalla_token", token)
	ctx = tflog.SetField(ctx, "sevalla_base_url", baseURL)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "sevalla_token")

	tflog.Debug(ctx, "Creating Sevalla client")

	// Create API client
	client := sevallaapi.NewClient(sevallaapi.Config{
		Token:   token,
		BaseURL: baseURL,
	})

	data_source_data := SevallaProviderData{
		Client: client,
	}

	resp.DataSourceData = data_source_data
	resp.ResourceData = data_source_data

	tflog.Info(ctx, "Configured Sevalla client", map[string]any{"success": true})
}

func (p *SevallaProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewApplicationResource,
		NewDatabaseResource,
		NewStaticSiteResource,
		NewObjectStorageResource,
		NewPipelineResource,
	}
}

func (p *SevallaProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewApplicationDataSource,
		NewDatabaseDataSource,
		NewStaticSiteDataSource,
		NewObjectStorageDataSource,
		NewPipelineDataSource,
	}
}

func (p *SevallaProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		// No functions for now
	}
}
