locals {
  cc_common = templatefile("${path.module}/cloud-config.yaml", {
    github_usernames = [for u in var.access_github_usernames : string(u) if u != ""]
    generic_bashrc   = filebase64("${path.module}/bashrc")
    generic_inputrc  = filebase64("${path.module}/inputrc")
    generic_aliases  = filebase64("${path.module}/bash_aliases")
  })
  merge_strategy = "list(append)+dict(no_replace,recurse_list)+str()"
}

