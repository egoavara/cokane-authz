
locals {
  app_name      = "cokane-authz"
  instance_name = var.name == null ? local.app_name : length(var.name) == 0 ? local.app_name : "${local.app_name}-${var.name}"
  version       = var.branch
}
