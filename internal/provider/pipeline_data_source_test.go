package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPipelineDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a pipeline first, then read it with data source
			{
				Config: testAccPipelineDataSourceConfig("test-pipeline-ds", "test-app-id-ds"),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check resource attributes
					resource.TestCheckResourceAttr("sevalla_pipeline.test", "name", "test-pipeline-ds"),
					resource.TestCheckResourceAttr("sevalla_pipeline.test", "app_id", "test-app-id-ds"),
					resource.TestCheckResourceAttr("sevalla_pipeline.test", "branch", "main"),
					resource.TestCheckResourceAttr("sevalla_pipeline.test", "auto_deploy", "true"),
					resource.TestCheckResourceAttrSet("sevalla_pipeline.test", "id"),
					resource.TestCheckResourceAttrSet("sevalla_pipeline.test", "created_at"),
					resource.TestCheckResourceAttrSet("sevalla_pipeline.test", "updated_at"),
					// Check data source attributes
					resource.TestCheckResourceAttr("data.sevalla_pipeline.test", "name", "test-pipeline-ds"),
					resource.TestCheckResourceAttr("data.sevalla_pipeline.test", "app_id", "test-app-id-ds"),
					resource.TestCheckResourceAttr("data.sevalla_pipeline.test", "branch", "main"),
					resource.TestCheckResourceAttr("data.sevalla_pipeline.test", "auto_deploy", "true"),
					resource.TestCheckResourceAttrSet("data.sevalla_pipeline.test", "id"),
					resource.TestCheckResourceAttrSet("data.sevalla_pipeline.test", "created_at"),
					resource.TestCheckResourceAttrSet("data.sevalla_pipeline.test", "updated_at"),
					// Check that resource and data source have the same ID
					resource.TestCheckResourceAttrPair("sevalla_pipeline.test", "id", "data.sevalla_pipeline.test", "id"),
				),
			},
		},
	})
}

func testAccPipelineDataSourceConfig(name, appID string) string {
	return providerConfig + fmt.Sprintf(`
resource "sevalla_pipeline" "test" {
  name        = %[1]q
  app_id      = %[2]q
  branch      = "main"
  auto_deploy = true
}

data "sevalla_pipeline" "test" {
  id = sevalla_pipeline.test.id
}
`, name, appID)
}
