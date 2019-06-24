output "PublicSubnetId" {
  value       = "${aws_subnet.public_subnet.id}"
  description = "ID of public subnet"
}
output "PrivateSubnetId" {
  value       = "${aws_subnet.private_subnet.id}"
  description = "ID of private subnet"
}
output "VpcId" {
  value       = "${aws_vpc.securus_vpc.id}"
  description = "ID of VPC"
}
output "PublicSubnetCidrBlock" {
  value       = "${aws_subnet.public_subnet.cidr_block}"
  description = "Public subnet cidr block"
}
output "PrivateSubnetCidrBlock" {
  value       = "${aws_subnet.private_subnet.cidr_block}"
  description = "Private subnet cidr block"
}
