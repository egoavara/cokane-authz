locals {
  app_label_selector = {
    "app.kubernetes.io/name"     = local.app_name
    "app.kubernetes.io/instance" = local.instance_name
    "app.kubernetes.io/version"  = local.version
  }
}

resource "kubernetes_deployment" "app" {
  metadata {
    namespace = var.namespace
    name      = "${local.instance_name}-deploy"
    labels = {
      app                           = local.instance_name
      "app.kubernetes.io/name"      = local.app_name
      "app.kubernetes.io/instance"  = local.instance_name
      "app.kubernetes.io/version"   = local.version
      "app.kubernetes.io/component" = "authorizer"
      "app.kubernetes.io/part-of"   = local.app_name
    }
  }
  spec {
    replicas = 1
    selector {
      match_labels = local.app_label_selector
    }
    template {
      metadata {
        labels = local.app_label_selector
      }
      spec {
        container {
          name              = "main"
          image             = var.image
          image_pull_policy = "Always"
          env {
            name  = "TEMPORAL_URL"
            value = var.temporal-url
          }
          port {
            name           = "http-api"
            container_port = 80
          }
          port {
            name           = "grpc-api"
            container_port = 50080
          }
        }
      }
    }
  }
}
