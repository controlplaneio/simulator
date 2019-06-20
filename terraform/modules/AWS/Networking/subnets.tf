resource "aws_subnet" "public_subnet" {
   vpc_id            = "${aws_vpc.securus_vpc.id}"
   cidr_block        = "${var.public_subnet_cidr}"
   availability_zone = "${var.public_avail_zone}"
   tags = {
     Name = "Securus Public subnet"   
   }
}

resource "aws_subnet" "private_subnet" {
   vpc_id            = "${aws_vpc.securus_vpc.id}"
   cidr_block        = "${var.private_subnet_cidr}"
   availability_zone = "${var.private_avail_zone}"
   tags = {
     Name = "Securus Private subnet"   
   }
}

