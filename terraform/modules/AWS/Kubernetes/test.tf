resource "null_resource" "master_test" {
  count = var.number_of_master_instances

  connection {
    bastion_host        = var.bastion_public_ip
    bastion_private_key = file(pathexpand("~/.kubesim/cp_simulator_rsa"))
    host = element(
      aws_instance.simulator_master_instances.*.private_ip,
      count.index,
    )
    type = "ssh"
    user = "root"

    // disable ssh-agent support
    agent       = "false"
    private_key = file(pathexpand("~/.kubesim/cp_simulator_rsa"))

    // Increase the timeout so the server has time to reboot
    timeout = "30m"
  }

  provisioner "file" {
    source      = "${path.module}/../../scripts/run-goss.sh"
    destination = "/root/run-goss.sh"
  }

  provisioner "file" {
    content = templatefile("${path.module}/master-goss.yaml", {
      "version_minor" = local.version_minor
    })
    destination = "/root/goss.yaml"
  }

  provisioner "remote-exec" {
    inline = [
      "set -o errexit",
      "chmod +x /root/run-goss.sh",
      "/root/run-goss.sh",
      "rm /root/run-goss.sh /root/goss.yaml",
    ]
  }
}

resource "null_resource" "node_test" {
  count = var.number_of_cluster_instances

  connection {
    bastion_host        = var.bastion_public_ip
    bastion_private_key = file(pathexpand("~/.kubesim/cp_simulator_rsa"))
    host = element(
      aws_instance.simulator_node_instances.*.private_ip,
      count.index,
    )
    type = "ssh"
    user = "root"

    // disable ssh-agent support
    agent       = "false"
    private_key = file(pathexpand("~/.kubesim/cp_simulator_rsa"))

    // Increase the timeout so the server has time to reboot
    timeout = "30m"
  }

  provisioner "file" {
    source      = "${path.module}/../../scripts/run-goss.sh"
    destination = "/root/run-goss.sh"
  }

  provisioner "file" {
    content = templatefile("${path.module}/node-goss.yaml", {
      "version_minor" = local.version_minor
    })
    destination = "/root/goss.yaml"
  }

  provisioner "remote-exec" {
    inline = [
      "set -o errexit",
      "chmod +x /root/run-goss.sh",
      "/root/run-goss.sh",
      "rm /root/run-goss.sh /root/goss.yaml",
    ]
  }
}

