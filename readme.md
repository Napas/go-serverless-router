A simple router written in Go for the [Serverless](https://serverless.com/) framework.

Currently due Go limitations needs to have a separate binary for each function, which is quite annoying. This library allows to have a single binary fo multiple functions.

## Installation
```go get github.com/Napas/go-serverless-router```

or if using [dep](https://github.com/golang/dep):
```dep ensure -v -add github.com/Napas/go-serverless-router```


## Supported Events
* APIGatewayProxyRequest
* DynamoDBEvent
* SQSEvent

Feel free to implement other if needed

## Bridges for the local development
Can be used with [LocalStack](https://github.com/localstack/localstack) for the local development
 
### Implemented bridges
* SQS

## Usage
```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	router "github.com/Napas/go-serverless-router"
	"github.com/Napas/go-serverless-router/routes"
	"github.com/Napas/go-serverless-router/bridges"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

const (
	envDev = "DEV"
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
	
	// CORS route for the /path/123
	corsRoute, err := routes.NewCorsApiGatewayRoute(
		"\\/path\\/\\d+", 
		"*", 
		[]string{"GET"}, 
		[]string{"Accept", "Content-Type"},
	)
	
	if err != nil {
		panic(err)
    }
    
	r.AddRoute(corsRoute)

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
	
	if os.Getenv("ENVIRONMENT") == envDev {
		// Local SQS client
		var sqsClient sqsiface.SQSAPI
		
		// This will consume messages from the queue and pass them to the router.
		// Should only be used for the local development
		sqsBridge := bridges.NewSqsBridge(
			r, 
			"http://localhost:3000/queue", 
			"arn:aws:sqs:us-east-2:123456789012:my-queue",
			sqsClient,
			"us-east-2",
			nil,
			)
		
		ctx := context.Background()
		sqsBridge.Run(ctx)
	}
	
	// Will match scheduled cloudwatch event with arn:aws:events:us-east-1:123456789012:rule/my-scheduled-rule
	cloudwatchScheduledEventRoute, err := routes.NewCloudwatchScheduledEventRoute(
		[]string{"^arn:aws:events:us-east-1:123456789012:rule\\/my-scheduled-rule$"},
		func(ctx context.Context, request events.CloudWatchEvent) error {
            // do something
            
            return nil
		},
	)
	
	if err != nil {
		panic(err)
	}
	
	r.AddRoute(cloudwatchScheduledEventRoute)

	// Start lambda with router as handler
	lambda.Start(r.Handle)
}

```
