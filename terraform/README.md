# Terraform environment  standup

**Table of Contents**

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


  - [Terraform target infrastructure](#terraform-target-infrastructure)
  - [Terraform Deployments structure](#terraform-deployments-structure)
  - [Terraform Module structure](#terraform-module-structure)
  - [Current Terraform Modules](#current-terraform-modules)
    - [SshKey](#sshkey)
    - [Ami](#ami)
    - [Bastion](#bastion)
    - [Kubernetes](#kubernetes)
    - [InternalHost](#internalhost)
    - [Networking](#networking)
    - [S3Storage](#s3storage)
    - [SecurityGroups](#securitygroups)
  - [Settings](#settings)
  - [Remote State](#remote-state)
- [Terraform environment  standup](#terraform-environment-standup)
  - [Terraform target infrastructure](#terraform-target-infrastructure)
  - [Terraform Deployments structure](#terraform-deployments-structure)
  - [Current Terraform Modules](#current-terraform-modules)
    - [SshKey](#sshkey)
    - [Ami](#ami)
    - [Bastion](#bastion)
    - [Kubernetes](#kubernetes)
    - [InternalHost](#internalhost)
    - [Networking](#networking)
    - [S3Storage](#s3storage)
    - [SecurityGroups](#securitygroups)
  - [Settings](#settings)
  - [Remote State](#remote-state)
- [Terraform Style Guide](#terraform-style-guide)
  - [Directory Structure](#directory-structure)
  - [Syntax](#syntax)
    - [Spacing](#spacing)
    - [Resource Block Alignment](#resource-block-alignment)
    - [Comments](#comments)
    - [Organizing Variables](#organizing-variables)
    - [Naming Conventions](#naming-conventions)
      - [File Names](#file-names)
      - [Parameter, Meta-parameter and User Variable Naming](#parameter-meta-parameter-and-user-variable-naming)
      - [Resource Naming](#resource-naming)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


## Terraform target infrastructure

![Terraform AWS infrastructure](./docs/aws-bastion-host-1.png)

## Terraform Deployments structure

Deployments are location within:

```deployments/[deployment name]```

The main.tf defines the modules that are to actioned and the variables passed to those modules.

## Terraform Module structure

Terraform modules are segregated by Cloud Provider under the modules directory, for example

```modules/AWS```

## Current Terraform Modules

The Terraform modules for AWS ( located under ```modules/AWS/[module name]``` ) action the following:


### SshKey

Refer to [settings documentation](./modules/AWS/SshKey/README-auto.md)

* Uploads Ssh key to be used for all instances

### Ami

This module does not require any settings to be provided at this moment

* Identifies current Ubuntu 18.04 LTS AMI Id for region in use

### Bastion

Refer to [settings documentation](./modules/AWS/Bastion/README-auto.md)

* A single bastion host on the public subnet

### Kubernetes

Refer to [settings documentation](./modules/AWS/Kubernetes/README-auto.md)


* One, or more, Kubernetes master nodes on the private network
* One, or more, Kubernetes nodes  on the private network

Cloud init is used to installed k8s software and initialise the cluster.  This is separted into 2 configurations:

* cloud-init-master.cfg - run on master nodes and installs kubelet, kubectl, kubeadm, docker and crictl.  Initialises the cluster
* cloud-init.cfg - runs on nodes and installs kubelet, kubectl, kubeadm, docker and crictl.

### InternalHost

Refer to [settings documentation](./modules/AWS/InternalHost/README-auto.md)

* A single host on the private subnet which is external to the cluster

### Networking

Refer to [settings documentation](./modules/AWS/Networking/README-auto.md)

* Single Vpc
* 2 subnets, 1 public, 1 private
* An Internet Gateway attached to public subnet
* A NAT gateway

The following routes are defined

* public_route_table - route to internet gateway, associated to public subnet
* private_nat_route_table - route to NAT gateway, associated to private subnet

### S3Storage

Refer to [settings documentation](./modules/AWS/S3Storage/README-auto.md)

* Create S3 bucket
* Create IAM role/policy for k8s hosts to access S3 bucket

### SecurityGroups

Refer to [settings documentation](./modules/AWS/SecurityGroups/README-auto.md)

The following security groups are defined

* bastion-sg - ingress connection over internet from defined cidr, open egress to internet
* controlplane-sg - Allows ingress from public subnet to private subnet (needs to be tighten up), open egress to internet (via NAT gateway)


## Settings

Refer to [settings documentation](./deployments/AWS/README-auto.md) for details on deployment settings required, and defaults provided.

## Remote State

This Terraform code uses remote state storage using S3. This is configured in `terraform/providers.tf` in the `terraform` block. The S3 bucket is assumed to have been created already.


# Terraform Style Guide

This following outlines coding conventions for Terraform's HashiCorp Configuration Language (HCL). Terraform allows infrastructure to be described as code. As such, we should adhere to a style guide to ensure readable and high quality code.

## Directory Structure

Directories are split into __deployments__ and __modules__

__deployments__ contains subdirectories for each terraform deployment

__modules__ contains subdirectories for each cloud vendor as appropriate (i.e AWS, Azure, GCP)

Taking __AWS__ as an example __modules__ subdirectory, the __AWS__ directory will then contains modules logically segregated into their core function (Networking, SecurityGroups etc).

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

### Naming Conventions

#### File Names

Create a separate resource file for each type of AWS resource. Similar resources should be defined in the same file.

```
main.tf
providers.tf
variables.tf
```

#### Parameter, Meta-parameter and User Variable Naming

 __Only use an underscore (`_`) when naming Terraform resources like TYPE/NAME parameters and user provided variables.__
 
```
resource "aws_security_group" "security_group" {
...
```

__Variables provided as output from modules should following CamelCase convention to make module provided variables easy to identify__

```
output "ControlPlaneSecurityGroupID" {
...
```

#### Resource Naming

__Only use a hyphen (`-`) when naming the component being created.__

```
resource "aws_security_group" "security_group" {
  name = "${var.resource_name}-security-group"
...
```

__A resource's NAME should describe TYPE pre-pended with simulator, unique ident and minus the provider.  Common shorthand for TYPE is fine here__

```
resource "aws_security_group" "simulator_controlplane_sg" {
...
```



