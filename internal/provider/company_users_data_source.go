package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sriniously/terraform-provider-sevalla/internal/sevallaapi"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CompanyUsersDataSource{}

func NewCompanyUsersDataSource() datasource.DataSource {
	return &CompanyUsersDataSource{}
}

// CompanyUsersDataSource defines the data source implementation.
type CompanyUsersDataSource struct {
	client *sevallaapi.Client
}

// CompanyUsersDataSourceModel describes the data source data model.
type CompanyUsersDataSourceModel struct {
	CompanyID types.String                 `tfsdk:"company_id"`
	Users     []CompanyUserDataSourceModel `tfsdk:"users"`
}

// CompanyUserDataSourceModel describes the user data model.
type CompanyUserDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	Email    types.String `tfsdk:"email"`
	Image    types.String `tfsdk:"image"`
	FullName types.String `tfsdk:"full_name"`
}

func (d *CompanyUsersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_company_users"
}

func (d *CompanyUsersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches the list of users for a company.",

		Attributes: map[string]schema.Attribute{
			"company_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier of the company.",
			},
			"users": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of users in the company.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier of the user.",
						},
						"email": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The email address of the user.",
						},
						"image": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The profile image URL of the user.",
						},
						"full_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The full name of the user.",
						},
					},
				},
			},
		},
	}
}

func (d *CompanyUsersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CompanyUsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CompanyUsersDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	users, err := d.client.Company.GetUsers(ctx, data.CompanyID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read company users, got error: %s", err))
		return
	}

	// Convert API users to terraform model
	var userModels []CompanyUserDataSourceModel
	for _, apiUser := range users.Company.Users {
		userModels = append(userModels, CompanyUserDataSourceModel{
			ID:       types.StringValue(apiUser.User.ID),
			Email:    types.StringValue(apiUser.User.Email),
			Image:    types.StringValue(apiUser.User.Image),
			FullName: types.StringValue(apiUser.User.FullName),
		})
	}

	data.Users = userModels

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
