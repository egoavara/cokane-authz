

resource "kubernetes_manifest" "manage-gateway" {
  manifest = {
    apiVersion = "networking.istio.io/v1alpha3",
    kind       = "Gateway",
    metadata = {
      namespace = kubernetes_namespace.ns.metadata[0].name,
      name      = "${local.instance_name}-manage-gateway",
    },
    spec = {
      selector = {
        "istio" = "ingressgateway"
      },
      servers = [
        {
          port = {
            number   = 80
            name     = "http"
            protocol = "HTTP"
          },
          hosts = [
            "auth.egoavara.net"
          ]
        }
      ]
    }
  }
}
