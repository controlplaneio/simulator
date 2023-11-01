resource "aws_vpc" "network" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true

  tags = merge(
    var.tags,
    {
      "Name" = title(var.name)
    }
  )
}

resource "aws_internet_gateway" "network" {
  vpc_id = aws_vpc.network.id

  tags = merge(
    var.tags,
    {
      "Name" = title(var.name)
    }
  )
}

resource "aws_eip" "network" {
  domain = "vpc"

  tags = merge(
    var.tags,
    {
      "Name" = title(var.name)
    }
  )

  depends_on = [
    aws_internet_gateway.network,
  ]
}

resource "aws_nat_gateway" "network" {
  subnet_id     = aws_subnet.public.id
  allocation_id = aws_eip.network.id

  tags = merge(
    var.tags,
    {
      "Name" = title(var.name)
    }
  )

  depends_on = [
    aws_eip.network,
  ]
}
