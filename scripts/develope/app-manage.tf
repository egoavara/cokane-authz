locals {
  app_label_selector = {
    "app.kubernetes.io/name"     = local.app_name
    "app.kubernetes.io/instance" = local.instance_name
    "app.kubernetes.io/version"  = local.version
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
        container {
          name    = "main"
          image   = "golang:1.23"
          command = ["/bin/sh"]
          args    = ["-c", file("${path.module}/files/entrypoint.sh")]
          volume_mount {
            name       = "git-sync-volume"
            mount_path = "/git"
          }
          port {
            protocol       = "TCP"
            name           = "http-app"
            container_port = 80
          }
        }
        container {
          name              = "git-sync"
          image             = "k8s.gcr.io/git-sync/git-sync:v4.2.4"
          image_pull_policy = "IfNotPresent"

          args = [
            "--repo=https://github.com/egoavara/cokane-authz",
            "--ref=refs/heads/${var.branch}",
            "--root=/git",
            "--period=10s",
            "--depth=1",
          ]
          volume_mount {
            name       = "git-sync-volume"
            mount_path = "/git"
          }
        }
        container {
          name  = "openfga"
          image = var.openfga-image
          args  = ["run"]
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
          env {
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
                name = "${local.app_name}-manage-datastore"
                key  = "USERNAME"
              }
            }
          }
          env {
            name = "OPENFGA_DATASTORE_PASSWORD"
            value_from {
              secret_key_ref {
                name = "${local.app_name}-manage-datastore"
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
                name = "${local.app_name}-manage-datastore"
                key  = "USERNAME"
              }
            }
          }

          env {
            name = "POSTGRES_PASSWORD"
            value_from {
              secret_key_ref {
                name = "${local.app_name}-manage-datastore"
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
            name = "${local.app_name}-manage-openfga-initdb"
          }
        }

        volume {
          name = "git-sync-volume"
          empty_dir {

          }
        }
      }
    }
  }
}
