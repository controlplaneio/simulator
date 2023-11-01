resource "aws_subnet" "private" {
  vpc_id            = aws_vpc.network.id
  cidr_block        = local.private_subnet_cidr
  availability_zone = var.availability_zone

  tags = merge(
    var.tags,
    {
      "Name" = format("%s Private", title(var.name))
    }
  )
}

resource "aws_route_table" "private" {
  vpc_id = aws_vpc.network.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.network.id
  }

  tags = merge(
    var.tags,
    {
      "Name" = format("%s Private", title(var.name))
    }
  )
}

resource "aws_route_table_association" "private" {
  route_table_id = aws_route_table.private.id
  subnet_id      = aws_subnet.private.id
}
