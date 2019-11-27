
locals {
  aws_tags = "${var.default_tags}"
}

// Setup networking
module "Networking" {
  source              = "../../modules/AWS/Networking"
  vpc_cidr            = "${var.vpc_cidr}"
  public_subnet_cidr  = "${var.public_subnet_cidr}"
  private_subnet_cidr = "${var.private_subnet_cidr}"
  default_tags        = "${local.aws_tags}"
}

// Discovery AMI Id to use for all instances
module "Ami" {
  source = "../../modules/AWS/Ami"
}

// Import ssh public key
module "SshKey" {
  source          = "../../modules/AWS/SshKey"
  access_key_name = "${var.access_key_name}"
  access_key      = "${var.access_key}"
}

// Setup Bastion host
module "Bastion" {
  source                   = "../../modules/AWS/Bastion"
  ami_id                   = "${module.Ami.AmiId}"
  instance_type            = "${var.instance_type}"
  access_key_name          = "${module.SshKey.KeyPairName}"
  security_group           = "${module.SecurityGroups.BastionSecurityGroupID}"
  subnet_id                = "${module.Networking.PublicSubnetId}"
  master_ip_addresses      = "${join(",", "${module.Kubernetes.K8sMasterPrivateIp}")}"
  node_ip_addresses        = "${join(",", "${module.Kubernetes.K8sNodesPrivateIp}")}"
  internal_node_private_ip = "${module.InternalNode.InternalNodePrivateIp}"
  attack_container_tag     = "${var.attack_container_tag}"
  default_tags             = "${local.aws_tags}"
}

// Setup Kubernetes master and nodes
module "Kubernetes" {
  source                      = "../../modules/AWS/Kubernetes"
  number_of_master_instances  = "${var.number_of_master_instances}"
  ami_id                      = "${module.Ami.AmiId}"
  master_instance_type        = "${var.master_instance_type}"
  number_of_cluster_instances = "${var.number_of_cluster_instances}"
  cluster_nodes_instance_type = "${var.cluster_nodes_instance_type}"
  bastion_public_ip           = "${module.Bastion.BastionPublicIp}"
  access_key_name             = "${module.SshKey.KeyPairName}"
  control_plane_sg_id         = "${module.SecurityGroups.ControlPlaneSecurityGroupID}"
  private_subnet_id           = "${module.Networking.PrivateSubnetId}"
  iam_instance_profile_id     = "${module.S3Storage.IamInstanceProfileId}"
  s3_bucket_name              = "${module.S3Storage.S3BucketName}"
  default_tags                = "${local.aws_tags}"
}

// Setup host within Kubernetes subnet
module "InternalNode" {
  source                   = "../../modules/AWS/InternalNode"
  ami_id                   = "${module.Ami.AmiId}"
  instance_type            = "${var.instance_type}"
  access_key_name          = "${module.SshKey.KeyPairName}"
  control_plane_sg_id      = "${module.SecurityGroups.ControlPlaneSecurityGroupID}"
  private_subnet_id        = "${module.Networking.PrivateSubnetId}"
  default_tags             = "${local.aws_tags}"
  bastion_public_ip        = "${module.Bastion.BastionPublicIp}"
  iam_instance_profile_id  = "${module.S3Storage.IamInstanceProfileId}"
  s3_bucket_name           = "${module.S3Storage.S3BucketName}"
}

// Create S3 bucket to share Kubernetes join details between
// master and nodes
module "S3Storage" {
  source         = "../../modules/AWS/S3Storage"
  default_tags   = "${local.aws_tags}"
}

// Define security groups
module "SecurityGroups" {
  source                    = "../../modules/AWS/SecurityGroups"
  access_cidr               = "${var.access_cidr}"
  vpc_id                    = "${module.Networking.VpcId}"
  public_subnet_cidr_block  = "${module.Networking.PublicSubnetCidrBlock}"
  private_subnet_cidr_block = "${module.Networking.PrivateSubnetCidrBlock}"
  default_tags              = "${local.aws_tags}"
}
