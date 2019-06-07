resource "aws_eip" "securus_eip" {
  vpc      = true
  depends_on = ["aws_internet_gateway.securus_igw"]
}
