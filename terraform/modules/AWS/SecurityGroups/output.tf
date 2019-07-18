output "BastionSecurityGroupID" {
  value       = "${aws_security_group.simulator_bastion_sg.id}"
  description = "Bastion security group id"
}
output "ControlPlaneSecurityGroupID" {
  value       = "${aws_security_group.simulator_controlplane_sg.id}"
  description = "Controlplane (Kubernetes) security group id"
}

