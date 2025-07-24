package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccApplicationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccApplicationResourceConfig("test-app"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sevalla_application.test", "display_name", "test-app"),
					resource.TestCheckResourceAttr("sevalla_application.test", "company_id", testAccCompanyID()),
					resource.TestCheckResourceAttr("sevalla_application.test", "repo_url", "https://github.com/test/test-app"),
					resource.TestCheckResourceAttr("sevalla_application.test", "auto_deploy", "true"),
					resource.TestCheckResourceAttrSet("sevalla_application.test", "name"),
					resource.TestCheckResourceAttrSet("sevalla_application.test", "id"),
					resource.TestCheckResourceAttrSet("sevalla_application.test", "status"),
					resource.TestCheckResourceAttrSet("sevalla_application.test", "created_at"),
					resource.TestCheckResourceAttrSet("sevalla_application.test", "updated_at"),
					// Check computed fields
					resource.TestCheckResourceAttrSet("sevalla_application.test", "default_branch"),
					resource.TestCheckResourceAttrSet("sevalla_application.test", "build_path"),
					resource.TestCheckResourceAttrSet("sevalla_application.test", "build_type"),
					resource.TestCheckResourceAttrSet("sevalla_application.test", "node_version"),
					resource.TestCheckResourceAttrSet("sevalla_application.test", "install_command"),
					// Check list attributes exist (even if empty)
					resource.TestCheckResourceAttrSet("sevalla_application.test", "environment_variables"),
					resource.TestCheckResourceAttrSet("sevalla_application.test", "deployments"),
					resource.TestCheckResourceAttrSet("sevalla_application.test", "processes"),
					resource.TestCheckResourceAttrSet("sevalla_application.test", "internal_connections"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sevalla_application.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccApplicationResourceConfig("test-app-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sevalla_application.test", "display_name", "test-app-updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccApplicationResourceConfig(name string) string {
	return providerConfig + fmt.Sprintf(`
resource "sevalla_application" "test" {
  display_name = %[1]q
  company_id   = %[2]q
  repo_url     = "https://github.com/test/test-app"
  auto_deploy  = true
}
`, name, testAccCompanyID())
}
