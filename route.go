package goserverlessrouter

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
)

var (
	validHttpMethods = []string{
		http.MethodGet,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPost,
		http.MethodConnect,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodTrace,
	}
)

type Route interface {
	Matches(request events.APIGatewayProxyRequest) bool
	GetHandler() HandlerFunc
}

type route struct {
	method  string
	path    *regexp.Regexp
	handler HandlerFunc
}

func isAllowedHttpMethod(httpMethodToCheck string) bool {
	for _, httpMethod := range validHttpMethods {
		if httpMethod == httpMethodToCheck {
			return true
		}
	}

	return false
}

func NewRoute(
	method string,
	path string,
	handler HandlerFunc,
) (Route, error) {
	reg, err := regexp.Compile(path)

	if err != nil {
		return nil, err
	}

	if !isAllowedHttpMethod(method) {
		return nil, fmt.Errorf(
			"not allowed http method \"%s\"",
			method,
		)
	}

	return &route{method, reg, handler}, nil
}

func (r *route) Matches(request events.APIGatewayProxyRequest) bool {
	return request.HTTPMethod == r.method && r.path.MatchString(request.Path)
}

func (r *route) GetHandler() HandlerFunc {
	return r.handler
}
