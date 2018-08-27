package sumologic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/nextgenhealthcare/sumologic-sdk-go"
)

func TestAccSumoLogicHostedCollector_basic(t *testing.T) {
	collectorName := fmt.Sprintf("collector-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSumoLogicHostedCollectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumoLogicHostedCollectorWithDefaults(collectorName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSumoLogicHostedCollectorExists("sumologic_hosted_collector.collector"),
				),
			},
		},
	})
}

func testAccCheckSumoLogicHostedCollectorExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*sumologic.Client)
		for _, r := range s.RootModule().Resources {
			i, _ := strconv.Atoi(r.Primary.ID)
			if _, _, err := client.GetHostedCollector(i); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckSumoLogicHostedCollectorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*sumologic.Client)

	for _, r := range s.RootModule().Resources {
		i, _ := strconv.Atoi(r.Primary.ID)
		if _, _, err := client.GetHostedCollector(i); err != nil {
			if err == sumologic.ErrCollectorNotFound {
				continue
			}
			return err
		}
		return fmt.Errorf("Collector still exists")
	}
	return nil
}

func testAccSumoLogicHostedCollectorWithDefaults(r string) string {
	return fmt.Sprintf(`
resource "sumologic_hosted_collector" "collector" {
  name = "%s"
}
`, r)
}
