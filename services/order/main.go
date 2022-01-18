package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	msk "github.com/aws/aws-sdk-go-v2/service/kafka"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go/aws/credentials"
	sigv4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/ddvkid/learngo/internal/util"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/aws_msk_iam"
	"os"
	"strings"
)

var (
	awsRegion = util.GetEnv("AWS_REGION", "us-east-2")
	topic     = "order"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion), config.WithEC2IMDSRegion())
	if err != nil {
		fmt.Printf("failed to load configuration, %v\n", err)
	}

	client := msk.NewFromConfig(cfg)
	clusterDetails, err := util.GetClusterConfig(client)
	if err != nil {
		fmt.Printf("failed to GetClusterConfig, %v\n", err)
	}
	fmt.Println(*clusterDetails.Brokers.BootstrapBrokerStringPublicSaslIam)

	var creds *credentials.Credentials
	if tokenFile, ok := os.LookupEnv("AWS_WEB_IDENTITY_TOKEN_FILE"); ok {
		stsClient := sts.NewFromConfig(cfg)
		arn := os.Getenv("AWS_ROLE_ARN")
		sessionName := "order"
		b, err := stscreds.IdentityTokenFile(tokenFile).GetIdentityToken()
		if err != nil {
			fmt.Println(err)
		}

		resp, err := stsClient.AssumeRoleWithWebIdentity(context.TODO(), &sts.AssumeRoleWithWebIdentityInput{
			RoleArn:          &arn,
			RoleSessionName:  &sessionName,
			WebIdentityToken: aws.String(string(b)),
		})
		if err != nil {
			fmt.Println(err)
		}
		creds = credentials.NewStaticCredentials(*resp.Credentials.AccessKeyId, *resp.Credentials.SecretAccessKey, *resp.Credentials.SessionToken)
	} else {
		creds = credentials.NewEnvCredentials()
	}

	dialer := &kafka.Dialer{
		SASLMechanism: &aws_msk_iam.Mechanism{
			Signer: sigv4.NewSigner(creds),
			Region: awsRegion,
		},
		TLS: &tls.Config{},
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: strings.Split(*clusterDetails.Brokers.BootstrapBrokerStringPublicSaslIam, ","),
		Topic:   topic,
		Dialer:  dialer,
	})

	fmt.Println("start consuming ... !!")

	defer func() {
		err := reader.Close()
		if err != nil {
			fmt.Println("Error closing consumer: ", err)
			return
		}
		fmt.Println("Consumer closed")
	}()

	for {
		m, err := reader.ReadMessage(context.TODO())
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}

}
