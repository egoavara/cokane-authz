variable "branch" {
  type     = string
  nullable = false
}

variable "host" {
  type     = string
  default = "auth.egoavara.net"
  nullable = false
  
}

variable "name" {
  type     = string
  default  = "cokane-authz"
  nullable = true
}

variable "namespace" {
  type     = string
  default  = "cokane-authz"
  nullable = false
}

variable "openfga-image" {
  type     = string
  default  = "openfga/openfga:v1.5.9"
  nullable = false
}