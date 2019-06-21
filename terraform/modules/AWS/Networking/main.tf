
##################################################
# VPC and Subnet creation
#

resource "aws_vpc" "securus_vpc" {
  cidr_block           = "${var.vpc_cidr}"
  enable_dns_support   = true
  enable_dns_hostnames = true
  tags = {
    Name = "Securus VPC"
  }
}

resource "aws_subnet" "public_subnet" {
   vpc_id            = "${aws_vpc.securus_vpc.id}"
   cidr_block        = "${var.public_subnet_cidr}"
   availability_zone = "${var.public_avail_zone}"
   tags = {
     Name = "Securus Public subnet"   
   }
}

resource "aws_subnet" "private_subnet" {
   vpc_id            = "${aws_vpc.securus_vpc.id}"
   cidr_block        = "${var.private_subnet_cidr}"
   availability_zone = "${var.private_avail_zone}"
   tags = {
     Name = "Securus Private subnet"   
   }
}

##################################################
# Elastic IP creation
#

resource "aws_eip" "securus_eip" {
  vpc      = true
  depends_on = ["aws_internet_gateway.securus_igw"]
}

##################################################
# Internet gateway
#

resource "aws_internet_gateway" "securus_igw" {
  vpc_id = "${aws_vpc.securus_vpc.id}"
  tags = {
        Name = "Securus InternetGateway"
    }
}

##################################################
# NAT gateway
#

resource "aws_nat_gateway" "securus_nat" {
    allocation_id = "${aws_eip.securus_eip.id}"
    subnet_id = "${aws_subnet.public_subnet.id}"
    depends_on = ["aws_internet_gateway.securus_igw"]
}

##################################################
# Route tables and associations
#

resource "aws_route_table" "public_route_table" {
  vpc_id = "${aws_vpc.securus_vpc.id}"
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.securus_igw.id}"
  }
  tags   = {
    Name = "Securus Public internet route table"
  }
}

resource "aws_route_table" "private_nat_route_table" {
  vpc_id = "${aws_vpc.securus_vpc.id}"
  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = "${aws_nat_gateway.securus_nat.id}"
  }
  tags   = {
    Name = "Securus Private NAT route table"
  }
}

# Associate public subnet to public route table
resource "aws_route_table_association" "public_rt_assoc" {
  subnet_id      = "${aws_subnet.public_subnet.id}"
  route_table_id = "${aws_route_table.public_route_table.id}"
}

# Associate private subnet to private route table
resource "aws_route_table_association" "private_rt_assoc" {
  subnet_id      = "${aws_subnet.private_subnet.id}"
  route_table_id = "${aws_route_table.private_nat_route_table.id}"
}

