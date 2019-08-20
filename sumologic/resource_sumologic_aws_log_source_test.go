package sumologic

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/nextgenhealthcare/sumologic-sdk-go"
)

func TestAccSumoLogicAWSLogSource_cloudtrail(t *testing.T) {

	if os.Getenv("AWS_S3_BUCKET") == "" {
		t.Skip("Environment Variable AWS_S3_BUCKET is not set")
	}
	awsS3Bucket := os.Getenv("AWS_S3_BUCKET")

	if os.Getenv("AWS_SUMOLOGIC_IAM_ROLE_ARN") == "" {
		t.Skip("Environment Variable AWS_SUMOLOGIC_IAM_ROLE_ARN is not set")
	}
	awsSumoLogicIAMRoleArn := os.Getenv("AWS_SUMOLOGIC_IAM_ROLE_ARN")

	collectorName := fmt.Sprintf("collector-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSumoLogicAWSLogSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumoLogicAWSLogSourceWithDefaults(collectorName, "AwsCloudTrailBucket", awsS3Bucket, awsSumoLogicIAMRoleArn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSumoLogicAWSLogSourceExists("sumologic_aws_log_source.source"),
				),
			},
		},
	})
}

func TestAccSumoLogicAWSLogSource_lb(t *testing.T) {

	if os.Getenv("AWS_S3_BUCKET") == "" {
		t.Skip("Environment Variable AWS_S3_BUCKET is not set")
	}
	awsS3Bucket := os.Getenv("AWS_S3_BUCKET")

	if os.Getenv("AWS_SUMOLOGIC_IAM_ROLE_ARN") == "" {
		t.Skip("Environment Variable AWS_SUMOLOGIC_IAM_ROLE_ARN is not set")
	}
	awsSumoLogicIAMRoleArn := os.Getenv("AWS_SUMOLOGIC_IAM_ROLE_ARN")

	collectorName := fmt.Sprintf("collector-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSumoLogicAWSLogSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumoLogicAWSLogSourceWithDefaults(collectorName, "AwsElbBucket", awsS3Bucket, awsSumoLogicIAMRoleArn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSumoLogicAWSLogSourceExists("sumologic_aws_log_source.source"),
				),
			},
		},
	})
}

func testAccCheckSumoLogicAWSLogSourceExists(n string) resource.TestCheckFunc {
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

		_, _, err := client.GetAWSLogSource(collectorID, sourceID)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckSumoLogicAWSLogSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*sumologic.Client)

	for _, r := range s.RootModule().Resources {
		i, _ := strconv.Atoi(r.Primary.ID)
		c, _ := strconv.Atoi(r.Primary.Attributes["collector_id"])
		if _, _, err := client.GetAWSLogSource(c, i); err != nil {
			if err == sumologic.ErrSourceNotFound {
				continue
			}
			return err
		}
		return fmt.Errorf("Source still exists")
	}
	return nil
}

func testAccSumoLogicAWSLogSourceWithDefaults(r string, sumoLogicLogContentType string, awsS3Bucket string, awsSumoLogicIAMRoleArn string) string {
	return fmt.Sprintf(`
resource "sumologic_hosted_collector" "collector" {
  name = "%[1]s"
}

resource "sumologic_aws_log_source" "source" {
  name = "test"
  collector_id = "${sumologic_hosted_collector.collector.id}"
  source_type = "Polling"
  scan_interval = 30000
  content_type = "%[2]s"
	cutoff_relative_time = "-0h"
  third_party_ref {
    resources {
      service_type = "%[2]s"
      path {
        type = "S3BucketPathExpression"
        bucket_name = "%[3]s"
        path_expression = "*"
      }
      authentication {
        type = "AWSRoleBasedAuthentication"
        role_arn = "%[4]s"
      }
    }
	}
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
`, r, sumoLogicLogContentType, awsS3Bucket, awsSumoLogicIAMRoleArn)
}
