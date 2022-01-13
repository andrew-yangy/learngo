locals {
  msk_cluster_name    = "learngo-kafka"
  kafka_version   = "2.6.2"
  client_subnets  = module.vpc.public_subnets
  number_of_nodes = length(module.vpc.public_subnets)
  volume_size     = 100
  instance_type   = "kafka.t3.small"
  encryption_in_transit_client_broker = "TLS"
  encryption_in_transit_in_cluster = true
}

data "aws_subnet" "this" {
  id =  module.vpc.public_subnets[0]
}

resource "aws_security_group" "this" {
  name_prefix = "${local.msk_cluster_name}-"
  vpc_id      = data.aws_subnet.this.vpc_id
}

resource "aws_security_group_rule" "msk-plain" {
  from_port         = 9092
  to_port           = 9092
  protocol          = "tcp"
  security_group_id = aws_security_group.this.id
  type              = "ingress"
  self              = true
}

resource "aws_security_group_rule" "msk-tls" {
  from_port         = 9094
  to_port           = 9094
  protocol          = "tcp"
  security_group_id = aws_security_group.this.id
  type              = "ingress"
  self              = true
}

resource "aws_security_group_rule" "zookeeper-plain" {
  from_port         = 2181
  to_port           = 2181
  protocol          = "tcp"
  security_group_id = aws_security_group.this.id
  type              = "ingress"
  self              = true
}

resource "aws_security_group_rule" "zookeeper-tls" {
  from_port         = 2182
  to_port           = 2182
  protocol          = "tcp"
  security_group_id = aws_security_group.this.id
  type              = "ingress"
  self              = true
}

resource "random_id" "configuration" {
  prefix      = "${local.msk_cluster_name}-"
  byte_length = 8

  keepers = {
    kafka_version     = local.kafka_version
  }
}

resource "aws_msk_configuration" "this" {
  kafka_versions    = [random_id.configuration.keepers.kafka_version]
  name              = random_id.configuration.dec

  lifecycle {
    create_before_destroy = true
  }
  server_properties = ""
}

resource "aws_msk_cluster" "this" {
  depends_on = [aws_msk_configuration.this]

  cluster_name           = local.msk_cluster_name
  kafka_version          = local.kafka_version
  number_of_broker_nodes = local.number_of_nodes

  broker_node_group_info {
    client_subnets  = local.client_subnets
    ebs_volume_size = local.volume_size
    instance_type   = local.instance_type
    security_groups = aws_security_group.this.*.id
  }

  configuration_info {
    arn      = aws_msk_configuration.this.arn
    revision = aws_msk_configuration.this.latest_revision
  }

  client_authentication {
    sasl {
      iam = true
      scram = false
    }
  }

  encryption_info {
    encryption_in_transit {
      client_broker = local.encryption_in_transit_client_broker
      in_cluster    = local.encryption_in_transit_in_cluster
    }
  }
}