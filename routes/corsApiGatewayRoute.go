package routes

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
	"strings"
)

type corsApiGatewayRoute struct {
	*ApiGatewayRoute
	origin  string
	methods string
	headers string
}

func NewCorsApiGatewayRoute(
	path string,
	origin string,
	methods []string,
	headers []string,
) (*ApiGatewayRoute, error) {
	route := &corsApiGatewayRoute{
		origin:  origin,
		methods: strings.Join(addOptionsMethod(methods), ", "),
		headers: strings.Join(headers, ", "),
	}

	apiGatewayRoute, err := NewApiGatewayRoute(
		path,
		http.MethodOptions,
		route.handler,
	)

	if err != nil {
		return nil, err
	}

	return apiGatewayRoute, nil
}

func (route *corsApiGatewayRoute) handler(
	ctx context.Context,
	request events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  route.origin,
			"Access-Control-Allow-Methods": route.methods,
			"Access-Control-Allow-Headers": route.headers,
		},
	}, nil
}

func addOptionsMethod(methods []string) []string {
	if !hasOptionsMethod(methods) {
		return append(methods, http.MethodOptions)
	}

	return methods
}

func hasOptionsMethod(methods []string) bool {
	for _, method := range methods {
		if strings.EqualFold(method, http.MethodOptions) {
			return true
		}
	}

	return false
}
