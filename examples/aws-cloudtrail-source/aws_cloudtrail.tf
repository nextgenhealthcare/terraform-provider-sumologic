resource "aws_s3_bucket" "security_logs" {
  bucket = "terraform-provider-sumologic-cloudtrail-${random_pet.name.id}"

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "aws:kms"
      }
    }
  }
}

data "aws_iam_policy_document" "security_logs" {
  statement {
    sid = "AclCheck"

    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["cloudtrail.amazonaws.com"]
    }

    actions = [
      "s3:GetBucketAcl",
    ]

    resources = [
      "${aws_s3_bucket.security_logs.arn}",
    ]
  }

  statement {
    sid = "Write"

    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["cloudtrail.amazonaws.com"]
    }

    actions = [
      "s3:PutObject",
    ]

    resources = [
      "${aws_s3_bucket.security_logs.arn}/AWSLogs/*",
    ]

    condition {
      test     = "StringEquals"
      variable = "s3:x-amz-acl"

      values = [
        "bucket-owner-full-control",
      ]
    }
  }
}

resource "aws_s3_bucket_policy" "security_logs" {
  bucket = "${aws_s3_bucket.security_logs.id}"
  policy = "${data.aws_iam_policy_document.security_logs.json}"
}

resource "aws_cloudtrail" "sumologic" {
  # AWS will error if the S3 bucket policy won't allow CloudTrail
  depends_on = ["aws_s3_bucket_policy.security_logs"]

  name                       = "sumologic"
  s3_bucket_name             = "${aws_s3_bucket.security_logs.id}"
  is_multi_region_trail      = true
  enable_log_file_validation = true
}
