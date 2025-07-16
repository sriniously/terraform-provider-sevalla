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
					resource.TestCheckResourceAttr("sevalla_application.test", "name", "test-app"),
					resource.TestCheckResourceAttr("sevalla_application.test", "description", "Test application"),
					resource.TestCheckResourceAttrSet("sevalla_application.test", "id"),
					resource.TestCheckResourceAttrSet("sevalla_application.test", "status"),
					resource.TestCheckResourceAttrSet("sevalla_application.test", "created_at"),
					resource.TestCheckResourceAttrSet("sevalla_application.test", "updated_at"),
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
					resource.TestCheckResourceAttr("sevalla_application.test", "name", "test-app-updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccApplicationResourceConfig(name string) string {
	return providerConfig + fmt.Sprintf(`
resource "sevalla_application" "test" {
  name        = %[1]q
  description = "Test application"
  instances   = 1
  memory      = 512
  cpu         = 250
}
`, name)
}

