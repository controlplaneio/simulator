output "bastion_public_ip" {
  value       = module.Bastion.BastionPublicIp
  description = "Bastion public IP"
}

output "master_nodes_private_ip" {
  value       = module.Networking.MasterIPAddresses
  description = "Master node private IP"
}

output "cluster_nodes_private_ip" {
  value       = module.Networking.NodeIPAddresses
  description = "Cluster node private IPs"
}

output "internal_host_private_ip" {
  value       = module.Networking.InternalIPAddress
  description = "Private Subnet node IP"
}

output "access_cidr" {
  value       = var.access_cidr
  description = "Remote access IP"
}

output "access_github_usernames" {
  value       = var.access_github_usernames
  description = "github usernames to be added to ssh authorized keys"
}

output "ami_id" {
  value       = module.Ami.AmiId
  description = "AMI used for all instances"
}

