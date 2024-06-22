provider "kubernetes" {
  config_path    = "~/.kube/config"
  config_context = "k3d-egoavara-net"
}

terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }
  backend "kubernetes" {
    config_path = "~/.kube/config"
    config_context = "k3d-egoavara-net"
    secret_suffix = "cokane-authnz"
    namespace = "terraform"
  }
}
