resource "kubernetes_manifest" "github-runner" {
  manifest = yamldecode(<<EOF
apiVersion: actions.summerwind.dev/v1alpha1
kind: RunnerDeployment
metadata:
 name: ${replace(var.github-repo, "/", "-")}-runner
 namespace: ${kubernetes_namespace.devenv.metadata[0].name}
spec:
 replicas: 2
 template:
   spec:
      labels:
        - kubernetes
      repository: ${var.github-repo}
     
EOF
)
}