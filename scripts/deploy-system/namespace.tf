resource "kubernetes_namespace" "system" {
  metadata {
    name = "cokane-system"
    labels = {
      "istio-injection" = "disabled"
    }
  }
}
