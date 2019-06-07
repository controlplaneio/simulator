resource "aws_internet_gateway" "igw" {
  vpc_id = "${aws_vpc.securus_vpc.id}"
  tags = {
        Name = "InternetGateway"
    }
}
