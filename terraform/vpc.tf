resource "aws_default_vpc" "default" {
  cidr_block = "${var.private_vpc_cidr}"
}

resource "aws_default_vpc" "public" {
  cidr_block = "${var.public_vpc_cidr}"
}
