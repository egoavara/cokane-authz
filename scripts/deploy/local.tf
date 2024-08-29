
locals {
  app_name      = "cokane-authz"
  instance_name = var.name == null ? local.app_name : length(var.name) == 0 ? local.app_name : "${local.app_name}-${var.name}"
  tag           = regex("(.+):(.+)", var.image)[1]
  semver        = try(regex("(?:v)?(([0-9]+)\\.([0-9]+)\\.([0-9]+)(?:-([0-9A-Za-z-]+(?:\\.[0-9A-Za-z-]+)*))?(?:\\+[0-9A-Za-z-]+))", local.tag)[0], null)
  version       = local.semver != null ? local.semver : local.tag
}
