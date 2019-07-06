module "Networking" {
  source                     = "../../modules/AWS/Networking"

  vpc_cidr                   = "${var.vpc_cidr}"
  public_subnet_cidr         = "${var.public_subnet_cidr}"
  private_subnet_cidr        = "${var.private_subnet_cidr}"
  public_avail_zone          = "${var.public_avail_zone}"
  private_avail_zone         = "${var.private_avail_zone}"
}

module "Bastion" {
  source                      = "../../modules/AWS/Bastion"
  ami_id                      = "${var.ami_id}"
  instance_type               = "${var.instance_type}"
  access_key_name             = "${var.access_key_name}"
  access_key                  = "${var.access_key}"
  security_group              = "${module.SecurityGroups.BastionSecurityGroupID}"
  subnet_id                   = "${module.Networking.PublicSubnetId}"
}

module "Kubernetes" {
  source                      = "../../modules/AWS/Kubernetes"
  region                      = "${var.region}"
  number_of_master_instances  = "${var.number_of_master_instances}"
  ami_id                      = "${var.ami_id}"
  master_instance_type        = "${var.master_instance_type}"
  number_of_cluster_instances = "${var.number_of_cluster_instances}"
  cluster_nodes_instance_type = "${var.cluster_nodes_instance_type}"
  key_pair_name               = "${module.Bastion.KeyPairName}"
  control_plane_sg_id         = "${module.SecurityGroups.ControlPlaneSecurityGroupID}"
  private_subnet_id           = "${module.Networking.PrivateSubnetId}"
  iam_instance_profile_id     = "${module.S3Storage.IamInstanceProfileId}"
  s3_bucket_name              = "${var.s3_bucket_name}"
}
module "S3Storage" {
  source                      = "../../modules/AWS/S3Storage"
  s3_bucket_name              = "${var.s3_bucket_name}"
}
module "SecurityGroups" {
  source                      = "../../modules/AWS/SecurityGroups"
  access_cidr                 = "${var.access_cidr}"
  vpc_id                      = "${module.Networking.VpcId}"
  public_subnet_cidr_block    = "${module.Networking.PublicSubnetCidrBlock}"
  private_subnet_cidr_block   = "${module.Networking.PrivateSubnetCidrBlock}"
}



