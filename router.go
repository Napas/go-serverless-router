package go_serverless_router

import (
	"context"
	"encoding/json"

	"github.com/Napas/go-serverless-router/routes"

	"github.com/joomcode/errorx"
)

var (
	RouterErrors = errorx.NewNamespace("router")

	RouterRouteNotFoundError = RouterErrors.NewType("route_not_found")
)

type Router interface {
	AddRoute(route routes.Route) Router
	Handle(ctx context.Context, event map[string]interface{}) (interface{}, error)
}

type router struct {
	routes []routes.Route
	logger Logger
}

func New() Router {
	return &router{logger: &NilLogger{}}
}

func NewWithLogger(logger Logger) Router {
	return &router{logger: logger}
}

func (router *router) AddRoute(route routes.Route) Router {
	router.routes = append(router.routes, route)

	return router
}

func (router *router) Handle(ctx context.Context, event map[string]interface{}) (interface{}, error) {
	encoded, _ := json.Marshal(event)
	router.logger.Printf("Got event: %s", encoded)

	for _, route := range router.routes {
		if route.Matches(event) {
			resp, err := route.Handle(ctx, event)

			if route.HasResponse() {
				return resp, err
			}

			return err, err
		}
	}

	router.logger.Println("Route was not found")

	err := RouterRouteNotFoundError.New("Route not found")

	// not sure at this point if AWS Lambda is expecting error as a first argument
	// or as a second.
	return err, err
}
