output "BastionSecurityGroupID"{
  value       = "${aws_security_group.bastion-sg.id}"
  description = "Bastion security group id"
}
output "ControlPlaneSecurityGroupID"{
  value       = "${aws_security_group.controlplane-sg.id}"
  description = "Controlplane (Kubernetes) security group id"
}

