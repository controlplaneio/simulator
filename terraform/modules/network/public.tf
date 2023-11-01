resource "aws_subnet" "public" {
  vpc_id            = aws_vpc.network.id
  cidr_block        = local.public_subnet_cidr
  availability_zone = var.availability_zone

  tags = merge(
    var.tags,
    {
      "Name" = format("%s Public", title(var.name))
    }
  )
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.network.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.network.id
  }

  tags = merge(
    var.tags,
    {
      "Name" = format("%s Public",title(var.name))
    }
  )
}

resource "aws_route_table_association" "public" {
  route_table_id = aws_route_table.public.id
  subnet_id      = aws_subnet.public.id
}
