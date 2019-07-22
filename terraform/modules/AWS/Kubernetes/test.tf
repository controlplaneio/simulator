resource "null_resource" "master_test" {
  count = "${var.number_of_master_instances}"

  connection {
    bastion_host = "${var.bastion_public_ip}"
    bastion_private_key = "${file(pathexpand("~/.ssh/cp_simulator_rsa"))}"
    host = "${element(aws_instance.simulator_master_instances.*.private_ip, count.index)}"
    type = "ssh"
    user = "root"
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
    source      = "${path.module}/master-goss.yaml"
    destination = "/root/goss.yaml"
  }

  provisioner "remote-exec" {
    inline = [
      "chmod +x /root/run-goss.sh",
      "/root/run-goss.sh",
    ]
  }
}
