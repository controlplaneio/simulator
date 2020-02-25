// generate random uuid to append to all names of resources
// to ensure unique

resource "random_uuid" "unique" {
}

// Public subnet (Bastion) security group
// restricts ingress to identifier ip address, egress open

resource "aws_security_group" "simulator_bastion_sg" {
  name   = "simulator-bastion-sg-${random_uuid.unique.result}"
  vpc_id = var.vpc_id

  tags = merge(
    var.default_tags,
    {
      "Name" = "Simulator Bastion Security Group"
    },
  )
}

resource "aws_security_group_rule" "allow_all_egress" {
  type        = "egress"
  from_port   = 0
  to_port     = 0
  protocol    = -1
  cidr_blocks = ["0.0.0.0/0"]

  security_group_id = aws_security_group.simulator_bastion_sg.id
}

resource "aws_security_group_rule" "allow_user_ip_ssh" {
  type        = "ingress"
  protocol    = "tcp"
  from_port   = 22
  to_port     = 22
  cidr_blocks = var.access_cidr

  security_group_id = aws_security_group.simulator_bastion_sg.id
}

// Private subnet security group
// Restricts ingress from public subnet using ssh
// Egress open (via NAT for internet)

resource "aws_security_group" "simulator_controlplane_sg" {
  name   = "simulator-controlplane-sg-${random_uuid.unique.result}"
  vpc_id = var.vpc_id

  ingress {
    protocol    = "tcp"
    from_port   = 22
    to_port     = 22
    cidr_blocks = [var.public_subnet_cidr_block]
  }

  ingress {
    protocol    = -1
    from_port   = 0
    to_port     = 0
    cidr_blocks = [var.private_subnet_cidr_block]
  }

  egress {
    protocol    = -1
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(
    var.default_tags,
    {
      "Name" = "simulator Controlplane Security Group"
    },
  )
}

