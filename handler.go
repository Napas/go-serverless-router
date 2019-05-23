package goserverlessrouter

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type HandlerFunc = func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

type Handler interface {
	Handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}
