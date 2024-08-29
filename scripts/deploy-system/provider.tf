provider "kubernetes" {
  config_path    = "~/.kube/config"
  config_context = "k3d-egoavara-net"
}

provider "helm" {
  kubernetes {
    config_path    = "~/.kube/config"
    config_context = "docker-desktop"
  }
}

provider "random" {}

terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.6"
    }
  }
  backend "kubernetes" {
    config_path    = "~/.kube/config"
    config_context = "k3d-egoavara-net"
    secret_suffix  = "cokane-authnz-system"
    namespace      = "terraform"
  }
}
