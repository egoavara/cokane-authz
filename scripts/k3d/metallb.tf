
resource "kubernetes_manifest" "ippool" {
  manifest = {
    apiVersion = "metallb.io/v1beta1"
    kind       = "IPAddressPool"
    metadata = {
      name      = "default"
      namespace = "metallb-system"
    }
    spec = {
      addresses = [
        cidrsubnet(var.ingress_cidr, 8, 255)
      ]
    }
  }
}
