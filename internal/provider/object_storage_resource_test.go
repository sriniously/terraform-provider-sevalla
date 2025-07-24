package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccObjectStorageResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccObjectStorageResourceConfig("test-bucket"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sevalla_object_storage.test", "name", "test-bucket"),
					resource.TestCheckResourceAttr("sevalla_object_storage.test", "region", "us-east-1"),
					resource.TestCheckResourceAttrSet("sevalla_object_storage.test", "id"),
					resource.TestCheckResourceAttrSet("sevalla_object_storage.test", "endpoint"),
					resource.TestCheckResourceAttrSet("sevalla_object_storage.test", "access_key"),
					resource.TestCheckResourceAttrSet("sevalla_object_storage.test", "secret_key"),
					resource.TestCheckResourceAttrSet("sevalla_object_storage.test", "created_at"),
					resource.TestCheckResourceAttrSet("sevalla_object_storage.test", "updated_at"),
					// Size and objects should be 0 for a new bucket
					resource.TestCheckResourceAttr("sevalla_object_storage.test", "size", "0"),
					resource.TestCheckResourceAttr("sevalla_object_storage.test", "objects", "0"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sevalla_object_storage.test",
				ImportState:       true,
				ImportStateVerify: true,
				// secret_key is sensitive, so we ignore it in import verification
				ImportStateVerifyIgnore: []string{"secret_key"},
			},
			// Update and Read testing
			{
				Config: testAccObjectStorageResourceConfig("test-bucket-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sevalla_object_storage.test", "name", "test-bucket-updated"),
					resource.TestCheckResourceAttr("sevalla_object_storage.test", "region", "us-east-1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccObjectStorageResourceMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create minimal object storage
			{
				Config: testAccObjectStorageResourceConfigMinimal("minimal-bucket"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sevalla_object_storage.test", "name", "minimal-bucket"),
					resource.TestCheckResourceAttrSet("sevalla_object_storage.test", "id"),
					resource.TestCheckResourceAttrSet("sevalla_object_storage.test", "endpoint"),
					resource.TestCheckResourceAttrSet("sevalla_object_storage.test", "access_key"),
					resource.TestCheckResourceAttrSet("sevalla_object_storage.test", "secret_key"),
					resource.TestCheckResourceAttrSet("sevalla_object_storage.test", "created_at"),
					resource.TestCheckResourceAttrSet("sevalla_object_storage.test", "updated_at"),
				),
			},
		},
	})
}

func TestAccObjectStorageResourceWithRegion(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create object storage with different region
			{
				Config: testAccObjectStorageResourceConfigWithRegion("region-bucket", "eu-west-1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sevalla_object_storage.test", "name", "region-bucket"),
					resource.TestCheckResourceAttr("sevalla_object_storage.test", "region", "eu-west-1"),
					resource.TestCheckResourceAttrSet("sevalla_object_storage.test", "id"),
					resource.TestCheckResourceAttrSet("sevalla_object_storage.test", "endpoint"),
					resource.TestCheckResourceAttrSet("sevalla_object_storage.test", "access_key"),
					resource.TestCheckResourceAttrSet("sevalla_object_storage.test", "secret_key"),
				),
			},
		},
	})
}

func testAccObjectStorageResourceConfig(name string) string {
	return providerConfig + fmt.Sprintf(`
resource "sevalla_object_storage" "test" {
  name   = %[1]q
  region = "us-east-1"
}
`, name)
}

func testAccObjectStorageResourceConfigMinimal(name string) string {
	return providerConfig + fmt.Sprintf(`
resource "sevalla_object_storage" "test" {
  name = %[1]q
}
`, name)
}

func testAccObjectStorageResourceConfigWithRegion(name, region string) string {
	return providerConfig + fmt.Sprintf(`
resource "sevalla_object_storage" "test" {
  name   = %[1]q
  region = %[2]q
}
`, name, region)
}
