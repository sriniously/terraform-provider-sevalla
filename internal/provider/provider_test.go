package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the Sevalla client is properly configured.
	providerConfig = `
provider "sevalla" {
  # token and base_url will be read from environment variables:
  # SEVALLA_TOKEN and SEVALLA_BASE_URL (if set)
}
`
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"sevalla": providerserver.NewProtocol6WithError(New("test")()),
}

func TestProvider(t *testing.T) {
	// Test that the provider can be instantiated
	provider := New("test")()
	if provider == nil {
		t.Fatal("Provider should not be nil")
	}
}

func testAccPreCheck(t *testing.T) {
	// Skip acceptance tests if SEVALLA_TOKEN is not set
	if os.Getenv("SEVALLA_TOKEN") == "" {
		t.Skip("SEVALLA_TOKEN environment variable must be set for acceptance tests")
	}
	// Skip acceptance tests if SEVALLA_COMPANY_ID is not set
	if os.Getenv("SEVALLA_COMPANY_ID") == "" {
		t.Skip("SEVALLA_COMPANY_ID environment variable must be set for acceptance tests")
	}
}

// testAccCompanyID returns the company ID from environment variables for testing
func testAccCompanyID() string {
	return os.Getenv("SEVALLA_COMPANY_ID")
}
