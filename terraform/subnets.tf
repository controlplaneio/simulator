resource "aws_subnet" "public_subnet" {
   vpc_id            = "${aws_vpc.public.id}"
   cidr_block        = "${var.public_subnet_cidr}"
   availability_zone = "${var.public_avail_zone}"
   tags = {
     Name = "Public subnet"   
   }
}

resource "aws_subnet" "private_subnet" {
   vpc_id            = "${aws_vpc.private.id}"
   cidr_block        = "${var.private_subnet_cidr}"
   availability_zone = "${var.private_avail_zone}"
   tags = {
     Name = "Private subnet"   
   }
}

