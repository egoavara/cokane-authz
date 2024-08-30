provider "kubernetes" {
  config_path    = "~/.kube/config"
  config_context = "k3d-egoavara-net"
}

provider "random" {}

terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.6.2"
    }
  }
  backend "kubernetes" {
    config_path    = "~/.kube/config"
    config_context = "k3d-egoavara-net"
    secret_suffix  = "cokane-authnz-develope"
    namespace      = "terraform"

  }
}
