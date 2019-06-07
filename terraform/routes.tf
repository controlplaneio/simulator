resource "aws_route" "internet_access" {
  route_table_id         = "${aws_vpc.public.main_route_table_id}"
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = "${aws_internet_gateway.igw.id}"
}

resource "aws_route_table" "private_route_table" {
    vpc_id = "${aws_vpc.private.id}"
    tags = {
        Name = "Private route table"
    }
}

resource "aws_route" "private_route" {
	route_table_id  = "${aws_route_table.private_route_table.id}"
	destination_cidr_block = "0.0.0.0/0"
	nat_gateway_id = "${aws_nat_gateway.nat.id}"
}

# Associate subnet public subnet to public route table
resource "aws_route_table_association" "public_subnet" {
    subnet_id = "${aws_subnet.public_subnet.id}"
    route_table_id = "${aws_vpc.public.main_route_table_id}"
}

# Associate subnet private subnet to private route table
resource "aws_route_table_association" "private" {
    subnet_id = "${aws_subnet.private_subnet.id}"
    route_table_id = "${aws_route_table.private_route_table.id}"
}

