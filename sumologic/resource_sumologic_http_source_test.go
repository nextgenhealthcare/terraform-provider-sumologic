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

func TestAccSumoLogicHTTPSource_basic(t *testing.T) {

	collectorName := fmt.Sprintf("collector-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSumoLogicHTTPSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumoLogicHTTPSourceWithDefaults(collectorName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSumoLogicHTTPSourceExists("sumologic_http_source.source"),
				),
			},
		},
	})
}

func testAccCheckSumoLogicHTTPSourceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Source ID is set")
		}

		client := testAccProvider.Meta().(*sumologic.Client)

		sourceID, _ := strconv.Atoi(rs.Primary.ID)
		collectorID, _ := strconv.Atoi(rs.Primary.Attributes["collector_id"])

		_, _, err := client.GetHTTPSource(collectorID, sourceID)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckSumoLogicHTTPSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*sumologic.Client)

	for _, r := range s.RootModule().Resources {
		i, _ := strconv.Atoi(r.Primary.ID)
		c, _ := strconv.Atoi(r.Primary.Attributes["collector_id"])
		if _, _, err := client.GetHTTPSource(c, i); err != nil {
			if err == sumologic.ErrSourceNotFound {
				continue
			}
			return err
		}
		return fmt.Errorf("Source still exists")
	}
	return nil
}

func testAccSumoLogicHTTPSourceWithDefaults(r string) string {
	return fmt.Sprintf(`
resource "sumologic_hosted_collector" "collector" {
  name = "%[1]s"
}

resource "sumologic_http_source" "source" {
  name = "test"
  collector_id = "${sumologic_hosted_collector.collector.id}"
	source_type = "HTTP"
	filter {
    filter_type = "Exclude"
    name = "No INFO"
    regexp = "(?s).*\\[INFO\\].*(?s)" 
  }

  filter {
    filter_type = "Exclude"
    name = "No DEBUG"
    regexp = "(?s).*\\[DEBUG\\].*(?s)" 
  }
}
`, r)
}
