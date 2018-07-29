data "aws_iam_policy_document" "sumologic_assume_role_policy" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "AWS"
      identifiers = ["926226587429"]
    }

    condition {
      test     = "StringEquals"
      variable = "sts:ExternalId"

      values = [
        "${var.sumologic_aws_external_id}",
      ]
    }
  }
}

resource "aws_iam_role" "sumologic" {
  name = "SumoLogicLogAccess2"

  assume_role_policy = "${data.aws_iam_policy_document.sumologic_assume_role_policy.json}"
}

data "aws_iam_policy_document" "sumologic" {
  statement {
    actions = [
      "s3:ListBucketVersions",
      "s3:ListBucket",
    ]

    resources = ["${aws_s3_bucket.security_logs.arn}"]
  }

  statement {
    actions = [
      "s3:GetObject",
      "s3:GetObjectVersion",
    ]

    resources = ["${aws_s3_bucket.security_logs.arn}/*"]
  }
}

resource "aws_iam_policy" "sumologic" {
  name        = "SumoLogicLogAccess"
  path        = "/"
  description = "Policy for SumoLogic accessing logs in S3"
  policy      = "${data.aws_iam_policy_document.sumologic.json}"
}

resource "aws_iam_role_policy_attachment" "sumologic" {
  role       = "${aws_iam_role.sumologic.name}"
  policy_arn = "${aws_iam_policy.sumologic.arn}"
}
