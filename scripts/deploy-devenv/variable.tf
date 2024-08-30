variable "namespace" {
  type     = string
  default  = "cokane-devenv"
  nullable = false
}

variable "github-repo" {
  type     = string
  default  = "egoavara/cokane-authz"
  nullable = false
  
}

variable "github-token" {
  type      = string
  sensitive = true
  nullable  = false
}
