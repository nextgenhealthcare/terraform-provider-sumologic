resource "sumologic_hosted_collector" "example" {
  name = "example-${random_pet.name.id}"
}

resource "sumologic_aws_log_source" "example" {
  # SumoLogic will error if the IAM policy isn't attached yet
  depends_on = ["aws_iam_role_policy_attachment.sumologic"]

  name                 = "CloudTrail"
  collector_id         = "${sumologic_hosted_collector.example.id}"
  category             = "cloudtrail/example"
  source_type          = "Polling"
  scan_interval        = -1
  content_type         = "AwsCloudTrailBucket"
  cutoff_relative_time = "-0h"

  third_party_ref {
    resources {
      service_type = "AwsCloudTrailBucket"

      path {
        type            = "S3BucketPathExpression"
        bucket_name     = "${aws_s3_bucket.security_logs.id}"
        path_expression = "AWSLogs/*/CloudTrail/*"
      }

      authentication {
        type     = "AWSRoleBasedAuthentication"
        role_arn = "${aws_iam_role.sumologic.arn}"
      }
    }
  }
}
