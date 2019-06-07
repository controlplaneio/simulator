resource "aws_vpc" "private" {
  cidr_block = "${var.private_vpc_cidr}"
  enable_dns_support = true
  enable_dns_hostnames = true
  tags = {
    Name = "PrivateVPC"
  }
}

resource "aws_vpc" "public" {
  cidr_block = "${var.public_vpc_cidr}"
  enable_dns_support = true
  enable_dns_hostnames = true
  tags = {
    Name = "PublicVPC"
  }
}
