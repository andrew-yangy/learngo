data "aws_route53_zone" "selected" {
  name         = "sporthub.tv"
  private_zone = false
}

data "aws_elb_hosted_zone_id" "main" {}

data "kubernetes_ingress" "learngo" {
  metadata {
    name = "learngo"
    namespace = "learngo"
  }
}

resource "aws_route53_record" "api_record" {
  zone_id = data.aws_route53_zone.selected.zone_id
  name    = "learngo.${data.aws_route53_zone.selected.name}"
  type    = "A"

  alias {
    name                   = kubernetes_ingress.learngo.status.0.load_balancer.0.ingress.0.hostname
    zone_id                = data.aws_elb_hosted_zone_id.main.id
    evaluate_target_health = true
  }
}