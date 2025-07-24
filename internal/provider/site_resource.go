package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sriniously/terraform-provider-sevalla/internal/sevallaapi"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SiteResource{}
var _ resource.ResourceWithImportState = &SiteResource{}

func NewSiteResource() resource.Resource {
	return &SiteResource{}
}

// SiteResource defines the resource implementation.
type SiteResource struct {
	client *sevallaapi.Client
}

// DomainModel represents a domain attached to an environment.
type DomainModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}

// EnvironmentModel represents a site environment.
type EnvironmentModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	DisplayName   types.String `tfsdk:"display_name"`
	IsPremium     types.Bool   `tfsdk:"is_premium"`
	IsBlocked     types.Bool   `tfsdk:"is_blocked"`
	Domains       types.List   `tfsdk:"domains"`
	PrimaryDomain types.Object `tfsdk:"primary_domain"`
}

// SiteResourceModel describes the resource data model.
type SiteResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	DisplayName  types.String `tfsdk:"display_name"`
	CompanyID    types.String `tfsdk:"company_id"`
	Status       types.String `tfsdk:"status"`
	Environments types.List   `tfsdk:"environments"`
}

func (r *SiteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site"
}

func (r *SiteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a WordPress site on Sevalla platform with full environment and domain management support.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the site.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique name of the site.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name of the site.",
			},
			"company_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The company ID that owns this site.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The current status of the site.",
			},
			"environments": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of environments for this WordPress site.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The environment ID.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The environment name.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The environment display name.",
						},
						"is_premium": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether this is a premium environment.",
						},
						"is_blocked": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether this environment is blocked.",
						},
						"domains": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "List of domains attached to this environment.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The domain ID.",
									},
									"name": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The domain name.",
									},
									"type": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The domain type.",
									},
								},
							},
						},
						"primary_domain": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "The primary domain for this environment.",
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The primary domain ID.",
								},
								"name": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The primary domain name.",
								},
								"type": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The primary domain type.",
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *SiteResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SiteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SiteResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := sevallaapi.CreateSiteRequest{
		CompanyID:   data.CompanyID.ValueString(),
		DisplayName: data.DisplayName.ValueString(),
	}

	tflog.Debug(ctx, "Creating site", map[string]interface{}{
		"company_id":   createReq.CompanyID,
		"display_name": createReq.DisplayName,
	})

	opResp, err := r.client.Sites.Create(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create site, got error: %s", err))
		return
	}

	// Wait for the operation to complete
	siteID, err := r.waitForOperation(ctx, opResp.OperationID)
	if err != nil {
		resp.Diagnostics.AddError("Operation Error", fmt.Sprintf("Site creation operation failed: %s", err))
		return
	}

	// Fetch the created site
	site, err := r.client.Sites.Get(ctx, siteID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read created site, got error: %s", err))
		return
	}

	// Map all fields from API response
	r.mapSiteToModel(ctx, &data, &site.Site)

	tflog.Trace(ctx, "Created site resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SiteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SiteResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	site, err := r.client.Sites.Get(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read site, got error: %s", err))
		return
	}

	// Map all fields from API response
	r.mapSiteToModel(ctx, &data, &site.Site)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SiteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SiteResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := sevallaapi.UpdateSiteRequest{
		DisplayName: stringPointer(data.DisplayName.ValueString()),
	}

	site, err := r.client.Sites.Update(ctx, data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update site, got error: %s", err))
		return
	}

	// Map all fields from API response
	r.mapSiteToModel(ctx, &data, &site.Site)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SiteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SiteResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Sites.Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete site, got error: %s", err))
		return
	}
}

func (r *SiteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// waitForOperation waits for an operation to complete and returns the resource ID
func (r *SiteResource) waitForOperation(ctx context.Context, operationID string) (string, error) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	timeout := time.After(10 * time.Minute)

	for {
		select {
		case <-ticker.C:
			op, err := r.client.Operations.GetStatus(ctx, operationID)
			if err != nil {
				return "", fmt.Errorf("failed to get operation status: %w", err)
			}

			switch op.Status {
			case "completed":
				// Extract site ID from operation data or resource_id
				if op.ResourceID != "" {
					return op.ResourceID, nil
				}
				// If ResourceID is not set, try to extract from Data
				if op.Data != nil {
					if dataMap, ok := op.Data.(map[string]interface{}); ok {
						if siteID, ok := dataMap["site_id"].(string); ok {
							return siteID, nil
						}
					}
				}
				return "", fmt.Errorf("operation completed but site ID not found")
			case "failed":
				if op.Error != nil {
					return "", fmt.Errorf("operation failed: %s", *op.Error)
				}
				return "", fmt.Errorf("operation failed with unknown error")
			}
		case <-timeout:
			return "", fmt.Errorf("operation timed out after 10 minutes")
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}

// mapSiteToModel maps API response to Terraform model
func (r *SiteResource) mapSiteToModel(ctx context.Context, data *SiteResourceModel, site *sevallaapi.SiteDetails) {
	data.ID = types.StringValue(site.ID)
	data.Name = types.StringValue(site.Name)
	data.DisplayName = types.StringValue(site.DisplayName)
	data.CompanyID = types.StringValue(site.CompanyID)
	data.Status = types.StringValue(site.Status)

	// Convert environments
	environments := make([]attr.Value, len(site.Environments))
	for i, env := range site.Environments {
		// Convert domains for this environment
		domains := make([]attr.Value, len(env.Domains))
		for j, domain := range env.Domains {
			domainObj, _ := types.ObjectValue(
				map[string]attr.Type{
					"id":   types.StringType,
					"name": types.StringType,
					"type": types.StringType,
				},
				map[string]attr.Value{
					"id":   types.StringValue(domain.ID),
					"name": types.StringValue(domain.Name),
					"type": types.StringValue(domain.Type),
				},
			)
			domains[j] = domainObj
		}
		domainsAttrTypes := map[string]attr.Type{
			"id":   types.StringType,
			"name": types.StringType,
			"type": types.StringType,
		}
		domainsList, _ := types.ListValue(types.ObjectType{AttrTypes: domainsAttrTypes}, domains)

		// Convert primary domain
		primaryDomainObj, _ := types.ObjectValue(
			map[string]attr.Type{
				"id":   types.StringType,
				"name": types.StringType,
				"type": types.StringType,
			},
			map[string]attr.Value{
				"id":   types.StringValue(env.PrimaryDomain.ID),
				"name": types.StringValue(env.PrimaryDomain.Name),
				"type": types.StringValue(env.PrimaryDomain.Type),
			},
		)

		// Create environment object
		envObj, _ := types.ObjectValue(
			map[string]attr.Type{
				"id":             types.StringType,
				"name":           types.StringType,
				"display_name":   types.StringType,
				"is_premium":     types.BoolType,
				"is_blocked":     types.BoolType,
				"domains":        types.ListType{ElemType: types.ObjectType{AttrTypes: domainsAttrTypes}},
				"primary_domain": types.ObjectType{AttrTypes: domainsAttrTypes},
			},
			map[string]attr.Value{
				"id":             types.StringValue(env.ID),
				"name":           types.StringValue(env.Name),
				"display_name":   types.StringValue(env.DisplayName),
				"is_premium":     types.BoolValue(env.IsPremium),
				"is_blocked":     types.BoolValue(env.IsBlocked),
				"domains":        domainsList,
				"primary_domain": primaryDomainObj,
			},
		)
		environments[i] = envObj
	}

	envAttrTypes := map[string]attr.Type{
		"id":           types.StringType,
		"name":         types.StringType,
		"display_name": types.StringType,
		"is_premium":   types.BoolType,
		"is_blocked":   types.BoolType,
		"domains": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"id":   types.StringType,
			"name": types.StringType,
			"type": types.StringType,
		}}},
		"primary_domain": types.ObjectType{AttrTypes: map[string]attr.Type{
			"id":   types.StringType,
			"name": types.StringType,
			"type": types.StringType,
		}},
	}
	data.Environments, _ = types.ListValue(types.ObjectType{AttrTypes: envAttrTypes}, environments)
}
