package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/cysp/adzerk-management-sdk-go/testserver"
)

func TestAdTypeResource(t *testing.T) {
	s := testserver.NewHttpTestServer()
	defer s.Close()

	name := "name"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testCombinedConfig(
					testProviderConfig(s.URL),
					testAdTypeResourceConfig(640, 480, nil),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("kevel_ad_type.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "kevel_ad_type.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testCombinedConfig(
					testProviderConfig(s.URL),
					testAdTypeResourceConfig(640, 480, &name),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("kevel_ad_type.test", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				Config: testCombinedConfig(
					testProviderConfig(s.URL),
					testAdTypeResourceConfig(640, 640, nil),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("kevel_ad_type.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAdTypeResourceConfig(width int32, height int32, name *string) string {
	widthField := fmt.Sprintf(`width = %d`, width)
	heightField := fmt.Sprintf(`height = %d`, height)
	nameField := ""
	if name != nil {
		nameField = fmt.Sprintf(`name = %q`, *name)
	}
	return testResourceConfig("ad_type", "test", widthField, heightField, nameField)
}
