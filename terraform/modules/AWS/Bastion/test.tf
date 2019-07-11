resource "null_resource" "bastion_test" {
  # Ensure we can SSH as root for the goss tests and also for preturb.sh
  connection {
    host    = "${aws_instance.simulator_bastion.public_ip}"
    type    = "ssh"
    user    = "root"
    agent   = "false"
    # Increase the timeout so the server has time to reboot
    timeout = "10m"
  }

  provisioner "file" {
    source      = "${ path.module }/../../scripts/run-goss.sh"
    destination = "/root/run-goss.sh"
  }

  provisioner "file" {
    source      = "${ path.module }/goss.yaml"
    destination = "/root/goss.yaml"
  }

  provisioner "remote-exec" {
    inline = [
      "chmod +x /root/run-goss.sh",
      "/root/run-goss.sh",
    ]
  }
}
