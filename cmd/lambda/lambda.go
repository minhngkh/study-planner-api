package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/rs/zerolog/log"

	"study-planner-api/internal/api"
	"study-planner-api/internal/database"
	handlerImpl "study-planner-api/internal/handler"
	_ "study-planner-api/internal/utils/env"
)

var httpLambda *httpadapter.HandlerAdapterV2

func init() {
	log.Printf("Echo cold start")

	database.Instance()

	impl := api.NewStrictHandler(handlerImpl.NewHandler(), nil)

	handler := api.NewEchoHandler()
	api.RegisterHandlers(handler, impl)

	httpLambda = httpadapter.NewV2(handler)
}

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return httpLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
