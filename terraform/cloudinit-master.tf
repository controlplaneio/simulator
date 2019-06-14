
data "template_file" "init-script-master" {
  template = "${file("cloud-init/cloud-init-master.cfg")}"
  vars = {
    REGION = "${var.region}"
  }
}


data "template_cloudinit_config" "cloudinit-securus-master" {

  gzip = false
  base64_encode = false
  
  part {
    filename     = "cloud-init.cfg"
    content_type = "text/cloud-config"
    content      = "${data.template_file.init-script-master.rendered}"
  }

}

