locals {
  manage-openfga-initdb = {
    name = "${local.app_name}-manage-openfga-initdb"
  }
  manage-datastore = {
    name = "${local.app_name}-manage-datastore"
  }

  app_label_selector = {
    "app.kubernetes.io/name"     = local.app_name
    "app.kubernetes.io/instance" = local.instance_name
    "app.kubernetes.io/version"  = local.version
  }
}

resource "random_password" "manage-datastore-username" {
  length  = 32
  special = false
}

resource "random_password" "manage-datastore-password" {
  length  = 128
  special = false
}

resource "kubernetes_secret" "manage-datastore" {
  metadata {
    namespace = var.namespace
    name      = local.manage-datastore.name
    labels = {
    }
  }
  data = {
    "USERNAME" = random_password.manage-datastore-username.result
    "PASSWORD" = random_password.manage-datastore-password.result
  }
  type = "Opaque"
}

resource "kubernetes_config_map" "manage-openfga-initdb" {
  metadata {
    namespace = var.namespace
    name      = local.manage-openfga-initdb.name
    labels = {
      app                           = local.instance_name
      "app.kubernetes.io/name"      = local.app_name
      "app.kubernetes.io/instance"  = local.instance_name
      "app.kubernetes.io/version"   = local.version
      "app.kubernetes.io/component" = "manage"
      "app.kubernetes.io/part-of"   = local.app_name
    }
  }
  data = {
    for filepath in fileset("${path.module}/files/openfga-schema", "*.sql") : filepath => file("${path.module}/files/openfga-schema/${filepath}")
  }
}
