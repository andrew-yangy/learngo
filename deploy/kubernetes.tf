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