variable "image" {
  type     = string
  default  = "ghcr.io/egoavara/cokane-authz:feature-1"
  nullable = false
}

variable "name" {
  type     = string
  default  = null
  nullable = true
}

variable "namespace" {
  type     = string
  default  = "cokane-authz"
  nullable = false
}

variable "temporal-url" {
  type     = string
  default  = "http://localhost:7233"
  nullable = false
}
 
