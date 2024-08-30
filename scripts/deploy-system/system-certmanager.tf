
resource "helm_release" "certmanager" {
  repository = "https://charts.jetstack.io"
  chart      = "cert-manager"
  name       = "cert-manager"
  version    = "1.14.4"

  namespace        = kubernetes_namespace.system.metadata[0].name
  create_namespace = false
  values           = [
    file("${path.module}/files/certmanager.yaml")
  ]