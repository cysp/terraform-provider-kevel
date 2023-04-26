package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/cysp/adzerk-management-sdk-go/testserver"
)

func TestUserResource(t *testing.T) {
	s := testserver.NewHttpTestServer()
	defer s.Close()

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testCombinedConfig(
					testProviderConfig(s.URL),
					testUserResourceConfig("one", "https://example.org/one"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("kevel_user.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "kevel_user.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testCombinedConfig(
					testProviderConfig(s.URL),
					testUserResourceConfig("two", "https://example.org/one/two"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testUserResourceConfig(email string, name string) string {
	emailField := fmt.Sprintf(`email = %q`, email)
	nameField := fmt.Sprintf(`name = %q`, name)
	return testResourceConfig("user", emailField, nameField)
}
