package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/cysp/adzerk-management-sdk-go/testserver"
)

func TestChannelSiteMapResource(t *testing.T) {
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
					testChannelSiteMapResourceConfig(1, 2, 5),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("kevel_channel_site_map.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "kevel_channel_site_map.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testCombinedConfig(
					testProviderConfig(s.URL),
					testChannelSiteMapResourceConfig(1, 2, 10),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("kevel_channel_site_map.test", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				Config: testCombinedConfig(
					testProviderConfig(s.URL),
					testChannelSiteMapResourceConfig(3, 2, 10),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("kevel_channel_site_map.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testChannelSiteMapResourceConfig(channelId int32, siteId int32, priority int32) string {
	channelIdField := fmt.Sprintf(`channel_id = %d`, channelId)
	siteIdField := fmt.Sprintf(`site_id = %d`, siteId)
	priorityField := fmt.Sprintf(`priority = %d`, priority)
	return testResourceConfig("channel_site_map", "test", channelIdField, siteIdField, priorityField)
}
