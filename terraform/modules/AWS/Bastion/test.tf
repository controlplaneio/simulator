resource "null_resource" "bastion_test" {
  triggers = {
    attack_container_tag = var.attack_container_tag
  }

  // Ensure we can SSH as root for the goss tests and also for preturb.sh
  // Ensure we can SSH as root for the goss tests and also for preturb.sh
  connection {
    host = aws_instance.simulator_bastion.public_ip
    type = "ssh"
    user = "root"

    // disable ssh-agent support
    agent       = false
    private_key = file(pathexpand("~/.kubesim/cp_simulator_rsa"))

    // Increase the timeout so the server has time to reboot
    timeout = "10m"
  }

  provisioner "file" {
    source      = "${path.module}/../../scripts/run-goss.sh"
    destination = "/root/run-goss.sh"
  }

  provisioner "file" {
    content = templatefile("${path.module}/goss.yaml", {
      attack_container_tag  = var.attack_container_tag
      attack_container_repo = var.attack_container_repo
    })
    destination = "/root/goss.yaml"
  }

  provisioner "remote-exec" {
    inline = [
      "chmod +x /root/run-goss.sh",
      "/root/run-goss.sh",
      "rm /root/run-goss.sh /root/goss.yaml",
    ]
  }
}

