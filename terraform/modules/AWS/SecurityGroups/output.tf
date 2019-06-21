output "BastionSecurityGroupID"{
  value = "${aws_security_group.bastion-sg.id}"
}
output "ControlPlaneSecurityGroupID"{
  value = "${aws_security_group.controlplane-sg.id}"
}

