resource "kubernetes_secret" "github" {
  metadata {
    name      = "controller-manager"
    namespace = kubernetes_namespace.devenv.metadata[0].name
  }

  data = {
    "github_token" = var.github-token
  }
}


resource "helm_release" "github-runner-controller" {
  repository = "https://actions-runner-controller.github.io/actions-runner-controller"
  chart      = "actions-runner-controller"
  name       = "actions-runner-controller"
  version    = "0.23.7"

  namespace        = kubernetes_namespace.devenv.metadata[0].name
  create_namespace = false
  values = [
    file("${path.module}/files/github-runner-controller.yaml")
    ]
}
