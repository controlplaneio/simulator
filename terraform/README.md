# Terraform AWS environment  standup

The Terraform templates stand up the following:

* Single Vpc
* 2 subnets, 1 public, 1 private
* An Internet Gateway attached to public subnet
* A NAT gateway
* A single bastion host on the public subnet
* One, or more, K8s master nodes on the private network
* One, or more, K8s nodes  on the private network 

The following security groups are defined

* bastion-sg - ingress connection over internet from defined cidr, open egress to internet
* controlplane-sg - Allows ingress from public subnet to private subnet (needs to be tighten up), open egress to internet (via NAT gateway)

The following routes are defined

* public_route_table - route to internet gateway, associated to public subnet
* private_nat_route_table - route to NAT gateway, associated to private subnet

Cloud init is used to installed k8s software and initialise the cluster.  This is separted into 2 configurations:

* cloud-init-master.cfg - run on master nodes and installs kubelet, kubectl, kubeadm, docker and crictl.  Initialises the cluster
* cloud-init.cfg - runs on nodes and installs kubelet, kubectl, kubeadm, docker and crictl.

The k8s software listed above needs to be revisited.

## Settings

Refer to [settings documentation](https://github.com/controlplaneio/simulator-standalone/blob/ansible/terraform/README-auto.md) for details on settings requires, and defaults provided.

