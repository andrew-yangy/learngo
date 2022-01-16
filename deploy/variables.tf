variable "region" {
  default     = "us-east-2"
  description = "AWS region"
}

variable "environment" {
  default = "stage"
  type    = string
}

variable "k8s_namespace" {
  type = string
}

variable "image_registry" {
  type = string
}

variable "aws_access_key_id" {
  type = string
  sensitive = true
}

variable "aws_secret_access_key" {
  type = string
  sensitive = true
}