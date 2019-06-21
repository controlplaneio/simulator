output "PublicSubnetId" {
  value = "${aws_subnet.public_subnet.id}"
}
output "PrivateSubnetId" {
  value = "${aws_subnet.private_subnet.id}"
}
output "VpcId" {
  value = "${aws_vpc.securus_vpc.id}"
}
output "PublicSubnetCidrBlock" {
  value = "${aws_subnet.public_subnet.cidr_block}"
}
output "PrivateSubnetCidrBlock" {
  value = "${aws_subnet.private_subnet.cidr_block}"
}
