locals {
  serviceAccount = "order"
  oidc_fully_qualified_subjects = format("system:serviceaccount:%s:%s", var.k8s_namespace, local.serviceAccount)
}

resource "aws_iam_role" "irsa" {
  name  = "${local.cluster_name}-irsa"
  assume_role_policy = jsonencode({
    Statement = [{
      Action = "sts:AssumeRoleWithWebIdentity"
      Effect = "Allow"
      Principal = {
        Federated = module.eks.oidc_provider_arn
      }
      Condition = {
        StringEquals = {
          format("%s:sub", replace(module.eks.cluster_oidc_issuer_url, "https://", "")) = local.oidc_fully_qualified_subjects
        }
      }
    }]
    Version = "2012-10-17"
  })
}

resource "aws_iam_policy" "kafka_cluster" {
  name        = "${local.cluster_name}-policy-kafka_cluster"
  description = "Policy for kafka cluster"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "kafka-cluster:*",
        ]
        Effect   = "Allow"
        Resource = "*"
      },
    ]
  })
}
resource "aws_iam_role_policy_attachment" "irsa" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonMSKFullAccess"
  role       = aws_iam_role.irsa.name
}

resource "aws_iam_role_policy_attachment" "kafka" {
  policy_arn = aws_iam_policy.kafka_cluster.arn
  role       = aws_iam_role.irsa.name
}