output "testing" {
  value = "testing"
}

output "bastion_public_ip" {
  value = "1.1.1.1"
}

output "cluster_nodes_private_ip" {
  value = [ "127.0.0.1", "127.0.0.2" ]
}

output "internal_node_private_ip" {
  value = "127.0.0.3"
}

output "master_nodes_private_ip" {
  value = [ "127.0.0.4" ]
}
