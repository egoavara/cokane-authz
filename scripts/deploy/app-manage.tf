locals {
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
    name      = "${local.instance_name}-manage-datastore"
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
    name      = "${local.instance_name}-manage-openfga-initdb"
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

resource "kubernetes_deployment" "manage" {
  metadata {
    namespace = var.namespace
    name      = "${local.instance_name}-manage"
    labels = {
      app                           = local.instance_name
      "app.kubernetes.io/name"      = local.app_name
      "app.kubernetes.io/instance"  = local.instance_name
      "app.kubernetes.io/version"   = local.version
      "app.kubernetes.io/component" = "manage"
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
        # container {
        #   name              = "main"
        #   image             = var.image
        #   image_pull_policy = "Always"
        #   env {
        #     name  = "TEMPORAL_URL"
        #     value = var.temporal-url
        #   }
        #   port {
        #     name           = "http-api"
        #     container_port = 80
        #   }
        #   port {
        #     name           = "grpc-api"
        #     container_port = 50080
        #   }
        # }
        container {
          name  = "openfga"
          image = var.openfga-image
          args = ["run"]
          # playground
          env {
            name  = "OPENFGA_PLAYGROUND_ENABLED"
            value = "false"
          }

          # http
          env {
            name  = "OPENFGA_HTTP_ADDR"
            value = "0.0.0.0:8080"
          }
          # grpc
          env{
            name  = "OPENFGA_GRPC_ADDR"
            value = "0.0.0.0:8081"
          }

          # pprof
          env {
            name  = "OPENFGA_PROFILER_ENABLED"
            value = "true"
          }

          env {
            name  = "OPENFGA_PROFILER_ADDR"
            value = "0.0.0.0:13001"
          }

          # datastore
          env {
            name  = "OPENFGA_DATASTORE_ENGINE"
            value = "postgres"
          }
          env {
            name  = "OPENFGA_DATASTORE_URI"
            value = "postgres://localhost:5432/postgres?sslmode=disable"
          }
          env {
            name = "OPENFGA_DATASTORE_USERNAME"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.manage-datastore.metadata[0].name
                key  = "USERNAME"
              }
            }
          }
          env {
            name = "OPENFGA_DATASTORE_PASSWORD"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.manage-datastore.metadata[0].name
                key  = "PASSWORD"
              }
            }
          }


          port {
            protocol       = "TCP"
            name           = "http-app"
            container_port = 8080
          }
          port {
            protocol       = "TCP"
            name           = "grpc-app"
            container_port = 8081
          }
          port {
            protocol       = "TCP"
            name           = "http-pprof"
            container_port = 13001
          }
        }

        container {
          name  = "openfga-datastore"
          image = "postgres:16"

          env {
            name  = "POSTGRES_DB"
            value = "postgres"
          }

          env {
            name = "POSTGRES_USER"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.manage-datastore.metadata[0].name
                key  = "USERNAME"
              }
            }
          }

          env {
            name = "POSTGRES_PASSWORD"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.manage-datastore.metadata[0].name
                key  = "PASSWORD"
              }
            }
          }
          volume_mount {
            name       = "manage-openfga-initdb"
            mount_path = "/docker-entrypoint-initdb.d"
          }

          port {
            protocol       = "TCP"
            name           = "tcp-datastore"
            container_port = 5432
          }
        }
        volume {
          name = "manage-openfga-initdb"
          config_map {
            name = kubernetes_config_map.manage-openfga-initdb.metadata[0].name
          }
        }
      }
    }
  }
}
