
resource "helm_release" "eck-operator" {
  depends_on = [kubernetes_namespace.operator-system]

  repository = "https://helm.elastic.co"
  chart      = "eck-operator"
  name       = "eck-operator"
  version    = "2.10.0"

  namespace        = kubernetes_namespace.system.metadata[0].name
  create_namespace = false
  values = [
    file("${path.module}/files/eck-operator.yaml")
  ]
}
