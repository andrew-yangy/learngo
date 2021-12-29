provider "kubernetes" {
  host                   = data.aws_eks_cluster.cluster.endpoint
  token                  = data.aws_eks_cluster_auth.cluster.token
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)
}

resource "kubernetes_namespace" "learngo" {
  metadata {
    labels = {
      name = "learngo"
      "istio-injection" = "enabled"
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
      version: "v1"
    }
  }

  spec {
    replicas = var.k8s_replicaCount

    selector {
      match_labels = {
        app = var.k8s_name
        version: "v1"
      }
    }

    template {
      metadata {
        labels = {
          app = var.k8s_name
          version: "v1"
        }
      }

      spec {
        container {
          image = var.k8s_image.repository
          name  = var.k8s_name
          image_pull_policy = "IfNotPresent"
          port {
            container_port = var.k8s_image.containerPort
          }
        }
        image_pull_secrets {
          name = "ecr-secret"
        }
      }
    }
  }
}

resource "kubernetes_service" "learngo" {
  metadata {
    namespace = var.k8s_namespace
    name = var.k8s_name
    labels = {
      app = var.k8s_name
      service = var.k8s_name
    }
  }
  spec {
    selector = {
      app = var.k8s_name
      version: "v1"
    }
    port {
      name = "http"
      port = 80
      target_port = var.k8s_image.containerPort
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

resource "helm_release" "istio" {
  name       = "learngo-istio"
  namespace = var.k8s_namespace
  chart      = "./kube"
  depends_on = [
    kubernetes_service.learngo,
  ]
}
