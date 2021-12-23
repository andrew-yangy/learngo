variable "region" {
  default     = "us-east-2"
  description = "AWS region"
}

variable "environment" {
  default = "stage"
  type = string
}

variable "k8s_namespace" {
  type = string
}

variable "k8s_name" {
  type = string
}

variable "k8s_replicaCount" {
  type = number
}

variable "k8s_image" {}
