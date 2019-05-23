package goserverlessrouter

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type Router interface {
	Handler
	AddRoute(route Route)
}

type router struct {
	routes []Route
}

func NewRouter() Router {
	return &router{}
}

func (r *router) AddRoute(route Route) {
	r.routes = append(r.routes, route)
}

func (r *router) Handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	for _, route := range r.routes {
		if route.Matches(request) {
			return route.GetHandler()(ctx, request)
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotFound,
	}, nil
}
