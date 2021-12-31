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

variable "image_registry" {
  type = string
}