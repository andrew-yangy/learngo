resource "aws_ecr_repository" "order" {
  name = "order"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}

#data "aws_ecr_authorization_token" "token" {
#}
#
#resource "kubernetes_secret" this {
#  metadata {
#    namespace = var.k8s_namespace
#    name = "ecr-secret"
#  }
#
#  data = {
#    ".dockerconfigjson" = jsonencode({
#      auths = {
#        (var.image_registry) = {
#          auth = data.aws_ecr_authorization_token.token.authorization_token
#        }
#      }
#    })
#  }
#
#  type = "kubernetes.io/dockerconfigjson"
#  lifecycle {
#    ignore_changes = [data]
#  }
#}