
# resource "helm_release" "jaeger-operator" {
#   depends_on = [kubernetes_namespace.operator-system]

#   repository = "https://jaegertracing.github.io/helm-charts"
#   chart      = "jaeger-operator"
#   name       = "jaeger-operator"
#   version    = "2.49.0"

#   namespace        = kubernetes_namespace.system.metadata[0].name
#   create_namespace = false
#   values = [
#     file("${path.module}/files/jaeger-operator.yaml")
#   ]
# }
