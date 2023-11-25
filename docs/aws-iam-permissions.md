# AWS IAM Permissions

In order for Simulator to create AMIs and to Terraform the required infrastructure, create an IAM Role with the
following permissions.

You can use the Terraform configuration [here](../terraform/workspaces/simulator-iam/main.tf) to do this and get the IAM
Role ARN from the output.

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "ec2:AllocateAddress",
        "ec2:AssociateRouteTable",
        "ec2:AttachInternetGateway",
        "ec2:AttachVolume",
        "ec2:AuthorizeSecurityGroupEgress",
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
        "ec2:TerminateInstances"
      ],
      "Effect": "Allow",
      "Resource": "*"
    },
    {
      "Action": "sts:GetCallerIdentity",
      "Effect": "Allow",
      "Resource": "*"
    },
    {
      "Action": [
        "s3:CreateBucket",
        "s3:DeleteObject",
        "s3:GetObject",
        "s3:ListAllMyBuckets",
        "s3:ListBucket",
        "s3:PutObject"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
```

## Using SSO

Terraform does not support SSO correctly until v1.6., you can however, export the required variables to make it work

Create an SSO profile as usual

```shell
[profile simulator-sso]
sso_start_url = ...
sso_region = ...
sso_account_id = ...
sso_role_name = ...
region = ...
output = ...
```

Then, before running a `simulator image` or `simulator infra` command, ensure you have the required environment
variables set by running the following.

```shell
export AWS_REGION=...
aws sso login --profile simulator-sso
source <(aws configure export-credentials --format env)
```

Alternatively, you can use an SSO profile to perform an STS Assume Role

```shell
[profile simulator]
role_arn = arn:aws:iam::<account-id>:role/simulator
source_profile = simulator-sso
```

This time, ensure you login to your SSO profile, and then generate the STS credentials

```shell
export AWS_REGION=...
aws sso login --profile simulator-sso
source <(aws configure export-credentials --profile simulator --format env)
```

