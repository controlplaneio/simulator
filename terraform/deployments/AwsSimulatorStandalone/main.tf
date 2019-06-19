module "Networking" {
  source = "../../modules/Networking"

  public_subnet_cidr         = "${var.public_subnet_cidr}"
  availability_zone          = "${var.availability_zone}"
  private_subnet_cidr        = "${var.private_subnet_cidr}"
  private_avail_zone         = "${var.private_avail_zone}"
  vpc_cidr                   = "${var.vpc_cidr}"
}

