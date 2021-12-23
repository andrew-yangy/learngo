provider "kubernetes" {
  host                   = data.aws_eks_cluster.cluster.endpoint
  token                  = data.aws_eks_cluster_auth.cluster.token
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)
}

resource "kubernetes_namespace" "learngo" {
  metadata {
    labels = {
      name = "learngo"
    }

    name = "learngo"
  }
}

resource "kubernetes_deployment" "learngo" {
  metadata {
    namespace = var.k8s_namespace
    name = var.k8s_name
    labels = {
      app = var.k8s_name
    }
  }

  spec {
    replicas = 2

    selector {
      match_labels = {
        name = var.k8s_name
      }
    }

    template {
      metadata {
        labels = {
          name = var.k8s_name
        }
      }

      spec {
        container {
          image = var.k8s_image.repository
          name  = "app"
          port {
            container_port = var.k8s_image.containerPort
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "learngo" {
  metadata {
    namespace = var.k8s_namespace
    name = var.k8s_name
  }
  spec {
    type = "NodePort"
    selector = {
      name = var.k8s_name
    }
    port {
      port        = 80
      target_port = var.k8s_image.containerPort
    }
  }
}

resource "kubernetes_ingress" "learngo" {
  metadata {
    namespace = var.k8s_namespace
    name = var.k8s_name
    annotations = {
      "alb.ingress.kubernetes.io/scheme" = "internet-facing"
      "kubernetes.io/ingress.class" = "alb"
    }
  }

  spec {
    rule {
      http {
        path {
          path = "/*"
          backend {
            service_name = var.k8s_name
            service_port = 80
          }
        }
      }
    }
  }
}

provider "helm" {
  kubernetes {
    host                   = data.aws_eks_cluster.cluster.endpoint
    cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)
    token                  = data.aws_eks_cluster_auth.cluster.token
  }
}

resource "helm_release" "alb-ingress" {
  name       = "alb-ingress"
  namespace  = var.k8s_namespace
  chart      = "aws-alb-ingress-controller"
  repository = "https://cloudnativeapp.github.io/charts/curated/"

  set {
    name  = "autoDiscoverAwsRegion"
    value = "true"
  }
  set {
    name  = "autoDiscoverAwsVpcID"
    value = "true"
  }
  set {
    name  = "clusterName"
    value = local.cluster_name
  }
}