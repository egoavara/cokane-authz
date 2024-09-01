
locals {
  app_name      = "cokane-authz"
  instance_name = replace(var.branch, "/", "-")
  version       = replace(var.branch, "/", "-")
}
