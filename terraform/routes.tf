
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

