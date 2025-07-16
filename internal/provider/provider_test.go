package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the Sevalla client is properly configured.
	providerConfig = `
provider "sevalla" {
  token = "test-token"
  base_url = "https://api.sevalla.com"
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