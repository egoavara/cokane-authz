resource "kubernetes_namespace" "system" {
  metadata {
    name = "cokane-system"
  }
}
