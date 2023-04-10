package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/cysp/adzerk-management-sdk-go/testserver"
)

func TestChannelResource(t *testing.T) {
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
					testChannelResourceConfig("one", nil),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("kevel_channel.test", "id"),
					resource.TestCheckResourceAttr("kevel_channel.test", "title", "one"),
					resource.TestCheckResourceAttrSet("kevel_channel.test", "ad_types.#"),
					resource.TestCheckResourceAttr("kevel_channel.test", "ad_types.#", "0"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "kevel_channel.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testCombinedConfig(
					testProviderConfig(s.URL),
					testChannelResourceConfig("two", []int32{}),
				),
			},
			{
				Config: testCombinedConfig(
					testProviderConfig(s.URL),
					testChannelResourceConfig("three", []int32{123, 234}),
				),
			},
			{
				Config: testCombinedConfig(
					testProviderConfig(s.URL),
					testChannelResourceConfig("four", []int32{234}),
				),
			},
			{
				Config: testCombinedConfig(
					testProviderConfig(s.URL),
					testChannelResourceConfig("five", []int32{}),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testChannelResourceConfig(title string, adTypes []int32) string {
	titleField := fmt.Sprintf(`title = %q`, title)
	adTypesField := fmt.Sprintf(`ad_types = [%s]`, strings.Join(Map(adTypes, func(v int32) string { return fmt.Sprintf("%d", v) }), ", "))
	return testResourceConfig("channel", titleField, adTypesField)
}
