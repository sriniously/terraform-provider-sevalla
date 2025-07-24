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
					resource.TestCheckResourceAttr("sevalla_database.test", "display_name", "test-db"),
					resource.TestCheckResourceAttr("sevalla_database.test", "company_id", testAccCompanyID()),
					resource.TestCheckResourceAttr("sevalla_database.test", "type", "postgresql"),
					resource.TestCheckResourceAttr("sevalla_database.test", "version", "14"),
					resource.TestCheckResourceAttr("sevalla_database.test", "location", "us-central1"),
					resource.TestCheckResourceAttr("sevalla_database.test", "resource_type", "db1"),
					resource.TestCheckResourceAttr("sevalla_database.test", "db_name", "testdb"),
					resource.TestCheckResourceAttr("sevalla_database.test", "db_user", "testuser"),
					resource.TestCheckResourceAttr("sevalla_database.test", "db_password", "test-password"),
					resource.TestCheckResourceAttrSet("sevalla_database.test", "id"),
					resource.TestCheckResourceAttrSet("sevalla_database.test", "name"),
					resource.TestCheckResourceAttrSet("sevalla_database.test", "internal_hostname"),
					resource.TestCheckResourceAttrSet("sevalla_database.test", "internal_port"),
					resource.TestCheckResourceAttrSet("sevalla_database.test", "external_hostname"),
					resource.TestCheckResourceAttrSet("sevalla_database.test", "external_port"),
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
				// db_password is not returned from API, so we ignore it in import
				ImportStateVerifyIgnore: []string{"db_password"},
			},
			// Update and Read testing
			{
				Config: testAccDatabaseResourceConfig("test-db-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sevalla_database.test", "display_name", "test-db-updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDatabaseResourceConfig(name string) string {
	return providerConfig + fmt.Sprintf(`
resource "sevalla_database" "test" {
  display_name    = %[1]q
  company_id      = %[2]q
  location        = "us-central1"
  resource_type   = "db1"
  type            = "postgresql"
  version         = "14"
  db_name         = "testdb"
  db_password     = "test-password"
  db_user         = "testuser"
}
`, name, testAccCompanyID())
}
