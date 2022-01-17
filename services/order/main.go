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
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/aws_msk_iam"
	"net/http"
	"os"
	"strings"
)

var (
	port      = "8080"
	awsRegion = util.GetEnv("AWS_REGION", "us-east-2")
	topic     = "abc"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello Andrew!!")
	})

	r.GET("/health", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(200, gin.H{"user": user, "value": value})
		} else {
			c.JSON(200, gin.H{"user": user, "status": "no value"})
		}
	})

	return r
}

func main() {
	fmt.Println("AWS_ROLE_ARN", os.Getenv("AWS_ROLE_ARN"))
	fmt.Println("AWS_WEB_IDENTITY_TOKEN_FILE", os.Getenv("AWS_WEB_IDENTITY_TOKEN_FILE"))
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion), config.WithEC2IMDSRegion())
	if err != nil {
		fmt.Printf("failed to load configuration, %v\n", err)
	}
	stsClient := sts.NewFromConfig(cfg)
	arn := os.Getenv("AWS_ROLE_ARN")
	sessionName := "order"
	b, err := stscreds.IdentityTokenFile(os.Getenv("AWS_WEB_IDENTITY_TOKEN_FILE")).GetIdentityToken()
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

	client := msk.NewFromConfig(cfg)
	clusterDetails, err := util.GetClusterConfig(client)
	if err != nil {
		fmt.Printf("failed to GetClusterConfig, %v\n", err)
	}
	fmt.Println(*clusterDetails.Brokers.BootstrapBrokerStringPublicSaslIam)
	sharedTransport := &kafka.Transport{
		SASL: &aws_msk_iam.Mechanism{
			Signer: sigv4.NewSigner(credentials.NewStaticCredentials(*resp.Credentials.AccessKeyId, *resp.Credentials.SecretAccessKey, *resp.Credentials.SessionToken)),
			Region: awsRegion,
		},
		TLS: &tls.Config{},
	}

	w := &kafka.Writer{
		Addr:      kafka.TCP(strings.Split(*clusterDetails.Brokers.BootstrapBrokerStringPublicSaslIam, ",")...),
		Topic:     topic,
		Balancer:  &kafka.LeastBytes{},
		Transport: sharedTransport,
	}

	err = w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("Key-A"),
			Value: []byte("Hello World!"),
		},
		kafka.Message{
			Key:   []byte("Key-B"),
			Value: []byte("One!"),
		},
		kafka.Message{
			Key:   []byte("Key-C"),
			Value: []byte("Two!"),
		},
	)
	if err != nil {
		fmt.Println("failed to write messages:", err)
	}

	if err := w.Close(); err != nil {
		fmt.Println("failed to close writer:", err)
	}

	dialer := &kafka.Dialer{
		SASLMechanism: &aws_msk_iam.Mechanism{
			Signer: sigv4.NewSigner(credentials.NewStaticCredentials(*resp.Credentials.AccessKeyId, *resp.Credentials.SecretAccessKey, *resp.Credentials.SessionToken)),
			Region: awsRegion,
		},
		TLS: &tls.Config{},
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: strings.Split(*clusterDetails.Brokers.BootstrapBrokerStringPublicSaslIam, ","),
		Topic:   topic,
		Dialer:  dialer,
	})
	defer reader.Close()

	fmt.Println("start consuming ... !!", reader.Config())
	//for {
	m, err := reader.ReadMessage(context.TODO())
	if err != nil {
		fmt.Println(err)
		//break
	}
	fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	//}
	//if err := reader.Close(); err != nil {
	//	fmt.Print("failed to close reader:", err)
	//}

	r := setupRouter()
	r.Run(":" + port)
	fmt.Printf("Starting Order service at: %s\n", port)
}
