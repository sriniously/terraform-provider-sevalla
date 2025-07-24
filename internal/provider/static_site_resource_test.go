package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStaticSiteResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccStaticSiteResourceConfig("test-site"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sevalla_static_site.test", "display_name", "test-site"),
					resource.TestCheckResourceAttr("sevalla_static_site.test", "company_id", testAccCompanyID()),
					resource.TestCheckResourceAttr("sevalla_static_site.test", "repo_url", "https://github.com/test/test-site"),
					resource.TestCheckResourceAttr("sevalla_static_site.test", "default_branch", "main"),
					resource.TestCheckResourceAttr("sevalla_static_site.test", "auto_deploy", "true"),
					resource.TestCheckResourceAttr("sevalla_static_site.test", "build_command", "npm run build"),
					resource.TestCheckResourceAttr("sevalla_static_site.test", "published_directory", "dist"),
					resource.TestCheckResourceAttr("sevalla_static_site.test", "node_version", "18.16.0"),
					resource.TestCheckResourceAttrSet("sevalla_static_site.test", "id"),
					resource.TestCheckResourceAttrSet("sevalla_static_site.test", "name"),
					resource.TestCheckResourceAttrSet("sevalla_static_site.test", "hostname"),
					resource.TestCheckResourceAttrSet("sevalla_static_site.test", "status"),
					resource.TestCheckResourceAttrSet("sevalla_static_site.test", "git_type"),
					resource.TestCheckResourceAttrSet("sevalla_static_site.test", "created_at"),
					resource.TestCheckResourceAttrSet("sevalla_static_site.test", "updated_at"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sevalla_static_site.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccStaticSiteResourceConfig("test-site-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sevalla_static_site.test", "display_name", "test-site-updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccStaticSiteResourceMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create minimal static site
			{
				Config: testAccStaticSiteResourceConfigMinimal("minimal-site"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sevalla_static_site.test", "display_name", "minimal-site"),
					resource.TestCheckResourceAttrSet("sevalla_static_site.test", "id"),
					resource.TestCheckResourceAttrSet("sevalla_static_site.test", "status"),
					resource.TestCheckResourceAttrSet("sevalla_static_site.test", "created_at"),
					resource.TestCheckResourceAttrSet("sevalla_static_site.test", "updated_at"),
				),
			},
		},
	})
}

func testAccStaticSiteResourceConfig(name string) string {
	return providerConfig + fmt.Sprintf(`
resource "sevalla_static_site" "test" {
  display_name        = %[1]q
  company_id          = %[2]q
  repo_url            = "https://github.com/test/test-site"
  default_branch      = "main"
  auto_deploy         = true
  build_command       = "npm run build"
  published_directory = "dist"
  node_version        = "18.16.0"
}
`, name, testAccCompanyID())
}

func testAccStaticSiteResourceConfigMinimal(name string) string {
	return providerConfig + fmt.Sprintf(`
resource "sevalla_static_site" "test" {
  display_name = %[1]q
  company_id   = %[2]q
  repo_url     = "https://github.com/test/minimal-site"
}
`, name, testAccCompanyID())
}
