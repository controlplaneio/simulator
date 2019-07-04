output "PublicSubnetId" {
  value       = "${aws_subnet.simulator_public_subnet.id}"
  description = "ID of public subnet"
}
output "PrivateSubnetId" {
  value       = "${aws_subnet.simulator_private_subnet.id}"
  description = "ID of private subnet"
}
output "VpcId" {
  value       = "${aws_vpc.simulator_vpc.id}"
  description = "ID of VPC"
}
output "PublicSubnetCidrBlock" {
  value       = "${aws_subnet.simulator_public_subnet.cidr_block}"
  description = "Public subnet cidr block"
}
output "PrivateSubnetCidrBlock" {
  value       = "${aws_subnet.simulator_private_subnet.cidr_block}"
  description = "Private subnet cidr block"
}
