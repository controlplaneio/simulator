output "network_id" {
  value = aws_vpc.network.id
}

output "network_cidr" {
  value = aws_vpc.network.cidr_block
}

output "public_subnet_id" {
  value = aws_subnet.public.id
}

output "private_subnet_id" {
  value = aws_subnet.private.id
}
