resource "aws_nat_gateway" "securus_nat" {
    allocation_id = "${aws_eip.securus_eip.id}"
    subnet_id = "${aws_subnet.public_subnet.id}"
    depends_on = ["aws_internet_gateway.securus_igw"]
}
