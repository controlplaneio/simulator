resource "null_resource" "internal_node_test" {

  connection {
    bastion_host         = "${aws_instance.simulator_bastion.public_ip}"
    bastion_private_key  = "${file(pathexpand("~/.ssh/cp_simulator_rsa"))}"
    host                 = "${aws_instance.simulator_internal_node.private_ip}"
    type                 = "ssh"
    user                 = "root"
    // disable ssh-agent support
    agent       = "false"
    private_key = "${file(pathexpand("~/.ssh/cp_simulator_rsa"))}"
    // Increase the timeout so the server has time to reboot
    timeout = "10m"
  }

  provisioner "file" {
    source      = "${path.module}/../../scripts/run-goss.sh"
    destination = "/root/run-goss.sh"
  }

  provisioner "file" {
    source      = "${path.module}/goss.yaml"
    destination = "/root/goss.yaml"
  }

  provisioner "remote-exec" {
    inline = [
      "chmod +x /root/run-goss.sh",
      "/root/run-goss.sh",
    ]
  }
}
