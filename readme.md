A simple router written in Go for the [Serverless](https://serverless.com/) framework.

Currently due Go limitations needs to have a separate binary for each function, which is quite annoying. This library allows to have a single binary fo multiple functions.

## Installation
```go get github.com/Napas/go-serverless-router```

or if using [dep](https://github.com/golang/dep):
```dep ensure -v -add github.com/Napas/go-serverless-router```


## Supported Events
* APIGatewayProxyRequest
* DynamoDBEvent

Feel free to implement other if needed

## Usage
```go
package main

import (
	"context"
	"fmt"
	"net/http"

	router "github.com/Napas/go-serverless-router"
	"github.com/Napas/go-serverless-router/routes"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	r := router.New()

	// Will match GET /path/123
	httpRoute, err := routes.NewApiGatewayRoute(
		"\\/path\\/\\d+",
		http.MethodGet,
		func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			return events.APIGatewayProxyResponse{
				Body: fmt.Sprintf("Got a request with id %d", request.PathParameters["id"]),
			}, nil
		},
	)

	if err != nil {
		panic(err)
	}

	r.AddRoute(httpRoute)

	// Will match events from account id 111111 and table called table-name
	dynamodbEventRoute, err := routes.NewDynamoDbRoute(
		"^arn:aws:dynamodb:us-east-1:111111:table\\/table-name\\/stream.*$",
		func(ctx context.Context, request events.DynamoDBEvent) {
			// do something

			// on fail instead returning error this handler needs to panic
			panic("Failed to consume event")
		},
	)

	if err != nil {
		panic(err)
	}

	r.AddRoute(dynamodbEventRoute)
	
	// Will match events for the arn:aws:sqs:us-east-2:123456789012:my-queue queue.
	sqsEventRoute, err := routes.NewSqsRoute(
		"^arn:aws:sqs:us-east-2:123456789012:my-queue$",
		func(ctx context.Context, request events.SQSEvent) error {
		    // do something
		    
		    return nil
		},
	)
	
	if err != nil {
		panic(err)
	}
	
	r.AddRoute(sqsEventRoute)

	// Start lambda with router as handler
	lambda.Start(r.Handle)
}

```