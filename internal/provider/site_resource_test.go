package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/cysp/adzerk-management-sdk-go/testserver"
)

func TestSiteResource(t *testing.T) {
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
					testSiteResourceConfig("one", "https://example.org/one"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("kevel_site.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "kevel_site.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testCombinedConfig(
					testProviderConfig(s.URL),
					testSiteResourceConfig("two", "https://example.org/one/two"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testSiteResourceConfig(title string, url string) string {
	titleField := fmt.Sprintf(`title = %q`, title)
	urlField := fmt.Sprintf(`url = %q`, url)
	return testResourceConfig("site", titleField, urlField)
}
