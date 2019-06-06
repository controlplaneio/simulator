resource "aws_vpc" "private" {
  cidr_block = "${var.private_vpc_cidr}"
}

resource "aws_vpc" "public" {
  cidr_block = "${var.public_vpc_cidr}"
}
