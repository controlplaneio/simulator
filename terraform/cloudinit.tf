
data "template_file" "init-script" {
  template = "${file("cloud-init/cloud-init.cfg")}"
  vars {
    REGION = "${var.region}"
  }
}


data "template_cloudinit_config" "cloudinit-securus" {

  gzip = false
  base64_encode = false
  
  part {
    filename     = "cloud-init.cfg"
    content_type = "text/cloud-config"
    content      = "${data.template_file.init-script.rendered}"
  }

}

