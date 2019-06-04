package routes

import (
	"context"

	"github.com/joomcode/errorx"
)

var (
	RouteErrors = errorx.NewNamespace("route")

	RouteCompileError   = RouteErrors.NewType("route_compile")
	RouteMarshalError   = RouteErrors.NewType("marshal")
	RouteUnmarshalError = RouteErrors.NewType("unmarshal")
)

type Route interface {
	Matches(event map[string]interface{}) bool
	Handle(ctx context.Context, event map[string]interface{}) (interface{}, error)
	HasResponse() bool
}
