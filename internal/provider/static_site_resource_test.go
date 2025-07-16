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
					resource.TestCheckResourceAttr("sevalla_static_site.test", "name", "test-site"),
					resource.TestCheckResourceAttr("sevalla_static_site.test", "branch", "main"),
					resource.TestCheckResourceAttr("sevalla_static_site.test", "build_dir", "dist"),
					resource.TestCheckResourceAttr("sevalla_static_site.test", "build_cmd", "npm run build"),
					resource.TestCheckResourceAttr("sevalla_static_site.test", "repository.url", "https://github.com/test/test-site"),
					resource.TestCheckResourceAttr("sevalla_static_site.test", "repository.type", "github"),
					resource.TestCheckResourceAttr("sevalla_static_site.test", "repository.branch", "main"),
					resource.TestCheckResourceAttrSet("sevalla_static_site.test", "id"),
					resource.TestCheckResourceAttrSet("sevalla_static_site.test", "domain"),
					resource.TestCheckResourceAttrSet("sevalla_static_site.test", "status"),
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
					resource.TestCheckResourceAttr("sevalla_static_site.test", "name", "test-site-updated"),
					resource.TestCheckResourceAttr("sevalla_static_site.test", "build_dir", "public"),
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
					resource.TestCheckResourceAttr("sevalla_static_site.test", "name", "minimal-site"),
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
  name      = %[1]q
  branch    = "main"
  build_dir = "dist"
  build_cmd = "npm run build"
  
  repository {
    url    = "https://github.com/test/test-site"
    type   = "github"
    branch = "main"
  }
}
`, name)
}

func testAccStaticSiteResourceConfigMinimal(name string) string {
	return providerConfig + fmt.Sprintf(`
resource "sevalla_static_site" "test" {
  name = %[1]q
}
`, name)
}

func testAccStaticSiteResourceConfigUpdated(name string) string {
	return providerConfig + fmt.Sprintf(`
resource "sevalla_static_site" "test" {
  name      = %[1]q
  branch    = "main"
  build_dir = "public"
  build_cmd = "npm run build"
  
  repository {
    url    = "https://github.com/test/test-site-updated"
    type   = "github"
    branch = "main"
  }
}
`, name)
}