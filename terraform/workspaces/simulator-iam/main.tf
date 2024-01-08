variable "name" {
  description = "The name used for the IAM Role and the IAM Policy."
  default     = "simulator"
}

resource "aws_iam_role" "simulator" {
  name               = var.name
  assume_role_policy = data.aws_iam_policy_document.simulator_assume_role.json
}

data "aws_iam_policy_document" "simulator_assume_role" {
  statement {
    actions = [
      "sts:AssumeRole",
    ]

    principals {
      type = "AWS"
      identifiers = [
        "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root",
      ]
    }
  }
}

data "aws_iam_policy_document" "simulator" {
  statement {
    actions = [
      "ec2:AllocateAddress",
      "ec2:AssociateRouteTable",
      "ec2:AttachInternetGateway",
      "ec2:AttachVolume",
      "ec2:AuthorizeSecurityGroupEgress",
      "ec2:AuthorizeSecurityGroupIngress",
      "ec2:AuthorizeSecurityGroupIngress",
      "ec2:CopyImage",
      "ec2:CreateImage",
      "ec2:CreateInternetGateway",
      "ec2:CreateKeyPair",
      "ec2:CreateNatGateway",
      "ec2:CreateRoute",
      "ec2:CreateRouteTable",
      "ec2:CreateSecurityGroup",
      "ec2:CreateSnapshot",
      "ec2:CreateSubnet",
      "ec2:CreateTags",
      "ec2:CreateVolume",
      "ec2:CreateVpc",
      "ec2:DeleteInternetGateway",
      "ec2:DeleteKeyPair",
      "ec2:DeleteNatGateway",
      "ec2:DeleteRouteTable",
      "ec2:DeleteSecurityGroup",
      "ec2:DeleteSnapshot",
      "ec2:DeleteSubnet",
      "ec2:DeleteVolume",
      "ec2:DeleteVpc",
      "ec2:DeregisterImage",
      "ec2:DescribeAddresses",
      "ec2:DescribeAvailabilityZones",
      "ec2:DescribeImageAttribute",
      "ec2:DescribeImages",
      "ec2:DescribeInstanceAttribute",
      "ec2:DescribeInstanceCreditSpecifications",
      "ec2:DescribeInstanceStatus",
      "ec2:DescribeInstanceTypes",
      "ec2:DescribeInstances",
      "ec2:DescribeInternetGateways",
      "ec2:DescribeKeyPairs",
      "ec2:DescribeNatGateways",
      "ec2:DescribeNetworkInterfaces",
      "ec2:DescribeRegions",
      "ec2:DescribeRouteTables",
      "ec2:DescribeSecurityGroupRules",
      "ec2:DescribeSecurityGroups",
      "ec2:DescribeSnapshots",
      "ec2:DescribeSubnets",
      "ec2:DescribeTags",
      "ec2:DescribeVolumes",
      "ec2:DescribeVpcAttribute",
      "ec2:DescribeVpcs",
      "ec2:DetachInternetGateway",
      "ec2:DetachVolume",
      "ec2:DisassociateAddress",
      "ec2:DisassociateRouteTable",
      "ec2:GetPasswordData",
      "ec2:ImportKeyPair",
      "ec2:ModifyImageAttribute",
      "ec2:ModifyInstanceAttribute",
      "ec2:ModifySnapshotAttribute",
      "ec2:ModifyVpcAttribute",
      "ec2:RegisterImage",
      "ec2:ReleaseAddress",
      "ec2:RevokeSecurityGroupEgress",
      "ec2:RevokeSecurityGroupIngress",
      "ec2:RunInstances",
      "ec2:StopInstances",
      "ec2:TerminateInstances",
    ]
    resources = [
      "*",
    ]
  }
  statement {
    actions = [
      "sts:GetCallerIdentity",
    ]
    resources = [
      "*",
    ]
  }
  statement {
    actions = [
      "s3:CreateBucket",
      "s3:DeleteBucket",
      "s3:DeleteObject",
      "s3:GetObject",
      "s3:ListAllMyBuckets",
      "s3:ListBucket",
      "s3:PutObject",
    ]
    resources = [
      "*",
    ]
  }
}

resource "aws_iam_policy" "simulator" {
  name        = var.name
  description = var.name
  policy      = data.aws_iam_policy_document.simulator.json
}

resource "aws_iam_role_policy_attachment" "simulator" {
  role       = aws_iam_role.simulator.name
  policy_arn = aws_iam_policy.simulator.arn
}

data "aws_caller_identity" "current" {}

output "simulator_iam_role_arn" {
  value = aws_iam_role.simulator.arn
}
