resource "kubernetes_service_account" this {
  metadata {
    namespace = var.k8s_namespace
    name = var.k8s_name
  }
  image_pull_secret {
    name = kubernetes_secret.this.metadata[0].name
  }
}

resource "kubernetes_deployment" this {
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
        service_account_name = kubernetes_service_account.this.metadata[0].name
        container {
          image = var.k8s_image.repository
          name  = var.k8s_name
          image_pull_policy = "IfNotPresent"
          port {
            container_port = var.k8s_image.containerPort
          }
        }
      }
    }
  }
}

resource "kubernetes_service" this {
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