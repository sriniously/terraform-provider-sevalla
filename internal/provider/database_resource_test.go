package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatabaseResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDatabaseResourceConfig("test-db"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sevalla_database.test", "name", "test-db"),
					resource.TestCheckResourceAttr("sevalla_database.test", "type", "postgresql"),
					resource.TestCheckResourceAttr("sevalla_database.test", "version", "14"),
					resource.TestCheckResourceAttrSet("sevalla_database.test", "id"),
					resource.TestCheckResourceAttrSet("sevalla_database.test", "host"),
					resource.TestCheckResourceAttrSet("sevalla_database.test", "port"),
					resource.TestCheckResourceAttrSet("sevalla_database.test", "username"),
					resource.TestCheckResourceAttrSet("sevalla_database.test", "status"),
					resource.TestCheckResourceAttrSet("sevalla_database.test", "created_at"),
					resource.TestCheckResourceAttrSet("sevalla_database.test", "updated_at"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sevalla_database.test",
				ImportState:       true,
				ImportStateVerify: true,
				// password is not returned from API, so we ignore it in import
				ImportStateVerifyIgnore: []string{"password"},
			},
			// Update and Read testing
			{
				Config: testAccDatabaseResourceConfig("test-db-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sevalla_database.test", "name", "test-db-updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDatabaseResourceConfig(name string) string {
	return providerConfig + fmt.Sprintf(`
resource "sevalla_database" "test" {
  name     = %[1]q
  type     = "postgresql"
  version  = "14"
  size     = "small"
  password = "test-password"
}
`, name)
}