package main

import (
	"context"
	"fmt"
	"github.com/ddvkid/learngo/internal/kafkaBuilder"
	"github.com/ddvkid/learngo/internal/util"
)

var (
	awsRegion = util.GetEnv("AWS_REGION", "us-east-2")
	topic     = "order"
)

func main() {
	handler, err := kafkaBuilder.NewKafkaHandler(kafkaBuilder.KafkaHandlerConfig{
		AwsRegion: awsRegion,
		Topic:     topic,
	})
	if err != nil {
		fmt.Println("Get kafka handler error: ", err)
	}

	handler.ReadMessage(context.Background())
}
