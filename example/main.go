package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"

	goserverlessrouter "github.com/Napas/go-serverless-router"
	"github.com/aws/aws-lambda-go/events"
)

func main() {
	// A route for the GET /
	// Second argument is a regex to match
	getRoute, err := goserverlessrouter.NewRoute(http.MethodGet, "\\/", GetHandler)

	if err != nil {
		panic(err)
	}

	// A route for the POST /
	postRoute, err := goserverlessrouter.NewRoute(http.MethodPost, "\\/", PostHandler)

	if err != nil {
		panic(err)
	}

	// A route for the GET /another
	anotherGetRoute, err := goserverlessrouter.NewRoute(http.MethodGet, "\\/another", AnotherGetHandler)

	if err != nil {
		panic(err)
	}

	router := goserverlessrouter.NewRouter()
	router.AddRoute(getRoute)
	router.AddRoute(postRoute)
	router.AddRoute(anotherGetRoute)

	lambda.Start(
		router.Handle,
	)
}

func GetHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Hello from GET handler",
	}, nil
}

func PostHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Hello from POST handler",
	}, nil
}

func AnotherGetHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Hello from Another GET handler",
	}, nil
}
