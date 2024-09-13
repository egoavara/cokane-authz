resource "kubernetes_namespace" "devenv" {
  metadata {
    name = "cokane-devenv"
    labels = {
      "istio-injection" = "disabled"
    }
  }
}
resource "kubernetes_namespace" "system" {
  metadata {
    name = "github-runner-system"
    labels = {
      "istio-injection" = "disabled"
    }
  }

}
