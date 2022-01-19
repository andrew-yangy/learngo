package kafkaBuilder

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	msk "github.com/aws/aws-sdk-go-v2/service/kafka"
	"github.com/aws/aws-sdk-go/aws/credentials"
	sigv4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/ddvkid/learngo/internal/aws"
	"github.com/ddvkid/learngo/internal/util"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/aws_msk_iam"
	"strings"
)

var awsRegion = util.GetEnv("AWS_REGION", "us-east-2")

type KafkaHandler struct {
	Reader *kafka.Reader
}

type KafkaHandlerConfig struct {
	AwsRegion string
	Topic     string
}

func NewKafkaHandler(handlerConfig KafkaHandlerConfig) (*KafkaHandler, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(handlerConfig.AwsRegion), config.WithEC2IMDSRegion())
	if err != nil {
		return nil, err
	}

	client := msk.NewFromConfig(cfg)
	clusterDetails, err := aws.GetClusterConfig(client)
	if err != nil {
		return nil, err
	}
	retrieveCreds, err := cfg.Credentials.Retrieve(context.Background())
	if err != nil {
		return nil, err
	}

	creds := credentials.NewStaticCredentials(retrieveCreds.AccessKeyID, retrieveCreds.SecretAccessKey, retrieveCreds.SessionToken)

	dialer := &kafka.Dialer{
		SASLMechanism: &aws_msk_iam.Mechanism{
			Signer: sigv4.NewSigner(creds),
			Region: awsRegion,
		},
		TLS: &tls.Config{},
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: strings.Split(*clusterDetails.Brokers.BootstrapBrokerStringPublicSaslIam, ","),
		Topic:   handlerConfig.Topic,
		Dialer:  dialer,
	})
	return &KafkaHandler{
		Reader: reader,
	}, nil
}

func (k *KafkaHandler) ReadMessage(ctx context.Context) {
	defer func() {
		err := k.Reader.Close()
		if err != nil {
			fmt.Println("Error closing consumer: ", err)
		}
		fmt.Println("Consumer closed")
	}()

	fmt.Println("start consuming ... !!")

	for {
		m, err := k.Reader.ReadMessage(ctx)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}
