resource "kubernetes_namespace" "devenv" {
  metadata {
    name = "cokane-devenv"
    labels = {
      "istio-injection" = "disabled"
    }
  }
}
