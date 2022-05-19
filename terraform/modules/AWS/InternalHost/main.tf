resource "aws_instance" "simulator_internal_host" {
  ami                         = var.ami_id
  key_name                    = var.access_key_name
  instance_type               = var.instance_type
  vpc_security_group_ids      = [var.control_plane_sg_id]
  associate_public_ip_address = false
  subnet_id                   = var.private_subnet_id
  user_data = templatefile("${path.module}/internal-config.yaml", {
    s3_bucket_name = var.s3_bucket_name
    host_bashrc    = filebase64("${path.module}/bashrc")
    host_inputrc   = filebase64("${path.module}/inputrc")
    host_aliases   = filebase64("${path.module}/bash_aliases")
  })
  iam_instance_profile = var.iam_instance_profile_id
  tags = merge(
    var.default_tags,
    {
      "Name" = "Simulator Internal Host"
    },
  )
}
