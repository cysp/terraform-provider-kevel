package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"kevel": providerserver.NewProtocol6WithError(New("test")()),
}

func testPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}

func testCombinedConfig(configs ...string) string {
	return strings.Join(configs, "\n\n")
}

func testProviderConfig(server string) string {
	return fmt.Sprintf(`
provider "kevel" {
	api_base_url = %[1]q
	api_key = "test"
}
`, server)
}

func testResourceConfig(resource string, name string, fields ...string) string {
	return fmt.Sprintf(`
resource "kevel_%s" "%s" {
	%s
}`, resource, name, strings.Join(fields, "\n  "))
}
