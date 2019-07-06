# Terraform environment  standup

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


### CreateBastion

Refer to [settings documentation](./modules/AWS/CreateBastion/README-auto.md)

* A single bastion host on the public subnet

### CreateK8s

Refer to [settings documentation](./modules/AWS/CreateK8s/README-auto.md)


* One, or more, K8s master nodes on the private network
* One, or more, K8s nodes  on the private network

Cloud init is used to installed k8s software and initialise the cluster.  This is separted into 2 configurations:

* cloud-init-master.cfg - run on master nodes and installs kubelet, kubectl, kubeadm, docker and crictl.  Initialises the cluster
* cloud-init.cfg - runs on nodes and installs kubelet, kubectl, kubeadm, docker and crictl.

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

This Terraform code uses remote state storage and locking using S3 and DynamoDB. This is configured in `terraform/providers.tf` in the `terraform` block. The S3 bucket and DynamoDB table are assumed to have been created already.

Terraform doesn't use the same AWS profile defined in `terraform/settings/bastion.tfvars`, it will use whatever is the default on your system. You may need to set the `AWS_PROFILE` environment variable to something different if you don't want it to use the default. - __This need to be confirmed as Terraform should work fine with the defined profile__

## Running the Terraform Code

To plan:

```bash
make infra-plan
```

To apply:
```bash
make infra-apply
```
To destroy:
```bash
make infra-destroy
```

