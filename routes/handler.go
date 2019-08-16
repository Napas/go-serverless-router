package routes

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type GeneralHandlerFunc func(ctx context.Context, request interface{}) (interface{}, error)
type ApiGatewayHandlerFunc func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
type DynamoDbHandlerFunc func(ctx context.Context, request events.DynamoDBEvent)
type SqsHandlerFunc func(ctx context.Context, request events.SQSEvent) error
type CloudWatchScheduledEventHandlerFunc func(ctx context.Context, request events.CloudWatchEvent) error

type GeneralHandler interface {
	Handle(ctx context.Context, request interface{}) (interface{}, error)
}

type ApiGatewayHandler interface {
	Handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

type DynamoDbHandler interface {
	Handle(ctx context.Context, request events.DynamoDBEvent)
}

type SqsHandler interface {
	Handle(ctx context.Context, request events.SQSEvent)
}

type CloudWatchScheduledEventHandler interface {
	Handle(ctx context.Context, request events.CloudWatchEvent) error
}
