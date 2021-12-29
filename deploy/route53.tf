data "aws_route53_zone" "selected" {
  name         = "sporthub.tv"
  private_zone = false
}

data "aws_elb_hosted_zone_id" "main" {}

data "kubernetes_service" "istio_ingress_gateway" {
  metadata {
    name        = "istio-ingressgateway"
    namespace   = "istio-system"
  }
  depends_on = [
    helm_release.istio,
  ]
}

resource "aws_route53_record" "api_record" {
  zone_id = data.aws_route53_zone.selected.zone_id
  name    = "learngo.${data.aws_route53_zone.selected.name}"
  type    = "A"

  alias {
    name                   = data.kubernetes_service.istio_ingress_gateway.status.0.load_balancer[0].ingress.0.hostname
    zone_id                = data.aws_elb_hosted_zone_id.main.id
    evaluate_target_health = true
  }
  depends_on = [
    helm_release.istio,
  ]
}