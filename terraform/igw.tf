resource "aws_internet_gateway" "securus_igw" {
  vpc_id = "${aws_vpc.securus_vpc.id}"
  tags = {
        Name = "Securus InternetGateway"
    }
}
