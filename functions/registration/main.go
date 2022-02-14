package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ddvkid/learngo/controllers/userregistration"
	"github.com/ddvkid/learngo/internal/storage"
	"github.com/ddvkid/learngo/internal/storage/postgres"
	"github.com/ddvkid/learngo/internal/util"
	"net/http"
)

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	store, err := postgres.NewPostgresStore()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, nil
	}
	return handle(ctx, req, store)
}

func handle(ctx context.Context, req events.APIGatewayProxyRequest, s storage.Store) (events.APIGatewayProxyResponse, error) {
	var registerReq userregistration.Request
	if err := util.DecodeRequest(req.Body, &registerReq); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}
	resp, err := userregistration.Register(ctx, s, registerReq)
	return util.EncodeLambdaResponse(resp, err)
}

func main() {
	lambda.Start(Handler)
}
