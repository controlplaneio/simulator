// identify availability zones

data "aws_availability_zones" "available" {
  state = "available"
}

// VPC and Subnet creation

resource "aws_vpc" "simulator_vpc" {
  cidr_block           = "${var.vpc_cidr}"
  enable_dns_support   = true
  enable_dns_hostnames = true
  tags = {
    Name = "Simulator VPC"
  }
}

resource "aws_subnet" "simulator_public_subnet" {
  vpc_id            = "${aws_vpc.simulator_vpc.id}"
  cidr_block        = "${var.public_subnet_cidr}"
  availability_zone = "${data.aws_availability_zones.available.names[0]}"
  tags = {
    Name = "Simulator Public subnet"
  }
}

resource "aws_subnet" "simulator_private_subnet" {
  vpc_id            = "${aws_vpc.simulator_vpc.id}"
  cidr_block        = "${var.private_subnet_cidr}"
  availability_zone = "${data.aws_availability_zones.available.names[1]}"
  tags = {
    Name = "Simulator Private subnet"
  }
}

// Elastic IP creation

resource "aws_eip" "simulator_eip" {
  vpc        = true
  depends_on = ["aws_internet_gateway.simulator_igw"]
}

// Internet gateway

resource "aws_internet_gateway" "simulator_igw" {
  vpc_id = "${aws_vpc.simulator_vpc.id}"
  tags = {
    Name = "Simulator InternetGateway"
  }
}

// NAT gateway

resource "aws_nat_gateway" "simulator_nat" {
  allocation_id = "${aws_eip.simulator_eip.id}"
  subnet_id     = "${aws_subnet.simulator_public_subnet.id}"
  depends_on    = ["aws_internet_gateway.simulator_igw"]
}

// Route tables and associations

resource "aws_route_table" "simulator_public_route_table" {
  vpc_id = "${aws_vpc.simulator_vpc.id}"
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.simulator_igw.id}"
  }
  tags = {
    Name = "Simulator Public internet route table"
  }
}

resource "aws_route_table" "simulator_private_nat_route_table" {
  vpc_id = "${aws_vpc.simulator_vpc.id}"
  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = "${aws_nat_gateway.simulator_nat.id}"
  }
  tags = {
    Name = "Simulator Private NAT route table"
  }
}

// Associate public subnet to public route table
resource "aws_route_table_association" "simulator_public_rt_assoc" {
  subnet_id      = "${aws_subnet.simulator_public_subnet.id}"
  route_table_id = "${aws_route_table.simulator_public_route_table.id}"
}

// Associate private subnet to private route table
resource "aws_route_table_association" "simulator_private_rt_assoc" {
  subnet_id      = "${aws_subnet.simulator_private_subnet.id}"
  route_table_id = "${aws_route_table.simulator_private_nat_route_table.id}"
}

