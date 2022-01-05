provider "kubernetes" {
  host                   = data.aws_eks_cluster.cluster.endpoint
  token                  = data.aws_eks_cluster_auth.cluster.token
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)
}

resource "kubernetes_namespace" "learngo" {
  metadata {
    labels = {
      name              = var.k8s_namespace
      "istio-injection" = "enabled"
    }

    name = var.k8s_namespace
  }
}

provider "helm" {
  kubernetes {
    host                   = data.aws_eks_cluster.cluster.endpoint
    cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)
    token                  = data.aws_eks_cluster_auth.cluster.token
  }
}

locals {
  chartHash = sha1(join("", [for f in fileset("./kube", "**/*.yaml") : filesha1("./kube/${f}")]))
}

resource "helm_release" "istio" {
  name      = "learngo-istio"
  namespace = var.k8s_namespace
  chart     = "./kube"

  set {
    name  = "chart hash"
    value = local.chartHash
  }
}

#module "order_service" {
#  source = "./modules/learngo-services"
#
#  app_name         = "order"
#  k8s_namespace    = var.k8s_namespace
#  k8s_name         = "order-api"
#  k8s_replicaCount = 2
#  k8s_image = {
#    repository    = "${var.image_registry}/order:latest"
#    containerPort = 8080
#  }
#  image_registry = var.image_registry
#}