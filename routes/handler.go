package routes

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type GeneralHandlerFunc func(ctx context.Context, request interface{}) (interface{}, error)
type ApiGatewayHandlerFunc func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
type DynamoDbHandlerFunc func(ctx context.Context, request events.DynamoDBEvent)

type GeneralHandler interface {
	Handle(ctx context.Context, request interface{}) (interface{}, error)
}

type ApiGatewayHandler interface {
	Handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

type DynamoDbHandler interface {
	Handle(ctx context.Context, request events.DynamoDBEvent)
}
