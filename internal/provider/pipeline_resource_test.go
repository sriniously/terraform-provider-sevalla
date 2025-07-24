package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPipelineResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccPipelineResourceConfig("test-pipeline", "test-app-id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sevalla_pipeline.test", "name", "test-pipeline"),
					resource.TestCheckResourceAttr("sevalla_pipeline.test", "app_id", "test-app-id"),
					resource.TestCheckResourceAttr("sevalla_pipeline.test", "branch", "main"),
					resource.TestCheckResourceAttr("sevalla_pipeline.test", "auto_deploy", "true"),
					resource.TestCheckResourceAttrSet("sevalla_pipeline.test", "id"),
					resource.TestCheckResourceAttrSet("sevalla_pipeline.test", "created_at"),
					resource.TestCheckResourceAttrSet("sevalla_pipeline.test", "updated_at"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sevalla_pipeline.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccPipelineResourceConfigUpdated("test-pipeline-updated", "test-app-id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sevalla_pipeline.test", "name", "test-pipeline-updated"),
					resource.TestCheckResourceAttr("sevalla_pipeline.test", "app_id", "test-app-id"),
					resource.TestCheckResourceAttr("sevalla_pipeline.test", "branch", "develop"),
					resource.TestCheckResourceAttr("sevalla_pipeline.test", "auto_deploy", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccPipelineResourceMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create minimal pipeline
			{
				Config: testAccPipelineResourceConfigMinimal("minimal-pipeline", "test-app-id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sevalla_pipeline.test", "name", "minimal-pipeline"),
					resource.TestCheckResourceAttr("sevalla_pipeline.test", "app_id", "test-app-id"),
					resource.TestCheckResourceAttr("sevalla_pipeline.test", "branch", "main"),
					resource.TestCheckResourceAttrSet("sevalla_pipeline.test", "id"),
					resource.TestCheckResourceAttrSet("sevalla_pipeline.test", "created_at"),
					resource.TestCheckResourceAttrSet("sevalla_pipeline.test", "updated_at"),
				),
			},
		},
	})
}

func TestAccPipelineResourceWithDifferentBranch(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create pipeline with different branch
			{
				Config: testAccPipelineResourceConfigWithBranch("branch-pipeline", "test-app-id", "feature/new-feature"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sevalla_pipeline.test", "name", "branch-pipeline"),
					resource.TestCheckResourceAttr("sevalla_pipeline.test", "app_id", "test-app-id"),
					resource.TestCheckResourceAttr("sevalla_pipeline.test", "branch", "feature/new-feature"),
					resource.TestCheckResourceAttr("sevalla_pipeline.test", "auto_deploy", "false"),
					resource.TestCheckResourceAttrSet("sevalla_pipeline.test", "id"),
				),
			},
		},
	})
}

func TestAccPipelineResourceAppIDReplacement(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create pipeline
			{
				Config: testAccPipelineResourceConfig("replacement-pipeline", "test-app-id-1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sevalla_pipeline.test", "name", "replacement-pipeline"),
					resource.TestCheckResourceAttr("sevalla_pipeline.test", "app_id", "test-app-id-1"),
				),
			},
			// Update app_id (should cause replacement)
			{
				Config: testAccPipelineResourceConfig("replacement-pipeline", "test-app-id-2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sevalla_pipeline.test", "name", "replacement-pipeline"),
					resource.TestCheckResourceAttr("sevalla_pipeline.test", "app_id", "test-app-id-2"),
				),
			},
		},
	})
}

func testAccPipelineResourceConfig(name, appID string) string {
	return providerConfig + fmt.Sprintf(`
resource "sevalla_pipeline" "test" {
  name        = %[1]q
  app_id      = %[2]q
  branch      = "main"
  auto_deploy = true
}
`, name, appID)
}

func testAccPipelineResourceConfigMinimal(name, appID string) string {
	return providerConfig + fmt.Sprintf(`
resource "sevalla_pipeline" "test" {
  name   = %[1]q
  app_id = %[2]q
  branch = "main"
}
`, name, appID)
}

func testAccPipelineResourceConfigUpdated(name, appID string) string {
	return providerConfig + fmt.Sprintf(`
resource "sevalla_pipeline" "test" {
  name        = %[1]q
  app_id      = %[2]q
  branch      = "develop"
  auto_deploy = false
}
`, name, appID)
}

func testAccPipelineResourceConfigWithBranch(name, appID, branch string) string {
	return providerConfig + fmt.Sprintf(`
resource "sevalla_pipeline" "test" {
  name        = %[1]q
  app_id      = %[2]q
  branch      = %[3]q
  auto_deploy = false
}
`, name, appID, branch)
}
