locals {
  public_subnet_cidr  = cidrsubnet(var.vpc_cidr, 1, 0)
  private_subnet_cidr = cidrsubnet(var.vpc_cidr, 1, 1)
}