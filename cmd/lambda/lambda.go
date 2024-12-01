package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/rs/zerolog/log"

	"study-planner-api/internal/server"
)

var httpLambda *httpadapter.HandlerAdapterV2

func init() {
	log.Printf("Echo cold start")
	e := server.NewEchoHandler()

	httpLambda = httpadapter.NewV2(e)
}

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return httpLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
