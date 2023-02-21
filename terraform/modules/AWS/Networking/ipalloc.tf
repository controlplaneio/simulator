locals {
  private_network_ip   = cidrhost(var.private_subnet_cidr, 0)
  private_network_mask = split("/", var.private_subnet_cidr)[1]
  max_hosts            = pow(2, 32 - local.private_network_mask)

  # First 3 and last 1 IPs are reserved
  # https://docs.aws.amazon.com/vpc/latest/userguide/configure-subnets.html
  private_network_min = 4
  private_network_max = local.max_hosts - 1

  total_hosts = (
    var.number_of_master_instances +
    var.number_of_cluster_instances +
    1
  )

  allocated_hosts = tolist([
    for i in random_shuffle.private_subnet_alloc.result : cidrhost(var.private_subnet_cidr, i)
  ])

  internal_ip = local.allocated_hosts[0]
  master_ips = slice(
    local.allocated_hosts,
    1,
    var.number_of_master_instances + 1
  )
  node_ips = slice(
    local.allocated_hosts,
    var.number_of_master_instances + 1,
    local.total_hosts
  )
}

resource "random_shuffle" "private_subnet_alloc" {
  seed         = aws_vpc.simulator_vpc.id
  input        = range(local.private_network_min, local.private_network_max)
  result_count = local.total_hosts
}
