# Terraform Style Guide

**Table of Contents**

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
- [Introduction](#introduction)
- [Syntax](#syntax)
  - [Spacing](#spacing)
  - [Resource Block Alignment](#resource-block-alignment)
  - [Comments](#comments)
  - [Organizing Variables](#organizing-variables)
- [Naming Conventions](#naming-conventions)
  - [File Names](#file-names)
  - [Parameter, Meta-parameter and Variable Naming](#parameter-meta-parameter-and-variable-naming)
  - [Resource Naming](#resource-naming)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Introduction

This outlines coding conventions for Terraform's HashiCorp Configuration Language (HCL). Terraform allows infrastructure to be described as code. As such, we should adhere to a style guide to ensure readable and high quality code.

## Syntax

- Strings are in double-quotes.

### Spacing

Use 2 spaces when defining resources except when defining inline policies or other inline resources.

```
resource "aws_iam_role" "iam_role" {
  name = "${var.resource_name}-role"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}
```

### Resource Block Alignment

Parameter definitions in a resource block should be aligned. The `terraform fmt` command can do this for you.

```
provider "aws" {
  access_key = "${var.aws_access_key}"
  secret_key = "${var.aws_secret_key}"
  region     = "us-east-1"
}
```


### Comments

When commenting use two "//" and a space in front of the comment.

```
// CREATE ELK IAM ROLE 
...
```

### Organizing Variables

The `variables.tf` file should be broken down into three sections with each section arranged alphabetically. Starting at the top of the file:

1. Variables that have no defaults defined
2. Variables that contain defaults
3. All locals blocks 

For example:

```
variable "image_tag" {}

variable "desired_count" {
  default = "2"
}

locals {
  domain_name = "${data.terraform_remote_state.account.domain_name}"
}
```

## Naming Conventions

### File Names

Create a separate resource file for each type of AWS resource. Similar resources should be defined in the same file.

```
main.tf
providers.tf
variables.tf
```

### Parameter, Meta-parameter and Variable Naming

 __Only use an underscore (`_`) when naming Terraform resources like TYPE/NAME parameters and variables.__
 
 ```
resource "aws_security_group" "security_group" {
...
```

### Resource Naming

__Only use a hyphen (`-`) when naming the component being created.__

 ```
resource "aws_security_group" "security_group" {
  name = "${var.resource_name}-security-group"
...
```

__A resource's NAME should be the same as the TYPE minus the provider.__

```
resource "aws_autoscaling_group" "autoscaling_group" {
...
```

If there are multiple resources of the same TYPE defined, add a minimalistic identifier to differentiate between the two resources. A blank line should sperate resource definitions contained in the same file.

```
// Create Data S3 Bucket
resource "aws_s3_bucket" "data_s3_bucket" {
  bucket = "${var.environment_name}-data-${var.aws_region}"
  acl    = "private"
  versioning {
    enabled = true
  }
}

// Create Images S3 Bucket
resource "aws_s3_bucket" "images_s3_bucket" {
  bucket = "${var.environment_name}-images-${var.aws_region}"
  acl    = "private"
}
```

