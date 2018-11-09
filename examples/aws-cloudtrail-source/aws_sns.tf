data "aws_caller_identity" "current" {}

resource "aws_sns_topic" "s3_updates" {
  name = "terraform-provider-sumologic-cloudtrail-${random_pet.name.id}"
}

data "aws_iam_policy_document" "s3_updates" {
  policy_id = "__default_policy_ID"

  statement {
    actions = [
      "SNS:Subscribe",
      "SNS:SetTopicAttributes",
      "SNS:RemovePermission",
      "SNS:Receive",
      "SNS:Publish",
      "SNS:ListSubscriptionsByTopic",
      "SNS:GetTopicAttributes",
      "SNS:DeleteTopic",
      "SNS:AddPermission",
    ]

    condition {
      test     = "StringEquals"
      variable = "AWS:SourceOwner"

      values = [
        "${data.aws_caller_identity.current.account_id}",
      ]
    }

    effect = "Allow"

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }

    resources = [
      "${aws_sns_topic.s3_updates.arn}",
    ]

    sid = "__default_statement_ID"
  }

  statement {
    actions = [
      "SNS:Publish",
    ]

    condition {
      test     = "ArnLike"
      variable = "aws:SourceArn"

      values = [
        "${aws_s3_bucket.security_logs.arn}",
      ]
    }

    effect = "Allow"

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }

    resources = [
      "${aws_sns_topic.s3_updates.arn}",
    ]

    sid = "AllowS3Updates"
  }
}

resource "aws_sns_topic_policy" "s3_updates" {
  arn = "${aws_sns_topic.s3_updates.arn}"

  policy = "${data.aws_iam_policy_document.s3_updates.json}"
}

resource "aws_sns_topic_subscription" "s3_updates" {
  topic_arn              = "${aws_sns_topic.s3_updates.arn}"
  protocol               = "https"
  endpoint               = "${sumologic_aws_log_source.example.url}"
  endpoint_auto_confirms = true
}
