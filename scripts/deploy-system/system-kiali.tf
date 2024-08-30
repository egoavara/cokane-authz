
# resource "helm_release" "kiali-operator" {
#   depends_on = [kubernetes_namespace.operator-system]

#   repository = "https://kiali.org/helm-charts"
#   chart      = "kiali-operator"
#   name       = "kiali-operator"
#   version    = "1.78.0"

#   namespace        = kubernetes_namespace.operator-system.metadata[0].name
#   create_namespace = false
#   values = [
#     file("${path.module}/files/kiali-operator.yaml")
#   ]
# }
