package goserverlessrouter_test

import (
	"context"
	"net/http"
	"reflect"
	"runtime"
	"testing"

	router "github.com/Napas/go-serverless-router"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func nilHandler(_ context.Context, _ events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	return response, err
}

func TestNewRoute(t *testing.T) {
	t.Run("Returns an error if invalid path regex is passed", func(t *testing.T) {
		invalidRegex := "["
		_, err := router.NewRoute(http.MethodGet, invalidRegex, nilHandler)

		assert.Error(t, err)
	})

	t.Run("Returns an error if invalid http method is passed", func(t *testing.T) {
		invalidHttpMethod := "TAKE"
		_, err := router.NewRoute(invalidHttpMethod, "\\/", nilHandler)

		assert.Error(t, err)
	})
}

func TestRoute_GetHandler(t *testing.T) {
	t.Run("Returns route handler", func(t *testing.T) {
		route, _ := router.NewRoute(http.MethodPost, "\\/", nilHandler)

		assert.Equal(
			t,
			runtime.FuncForPC(reflect.ValueOf(nilHandler).Pointer()).Name(),
			runtime.FuncForPC(reflect.ValueOf(route.GetHandler()).Pointer()).Name(),
		)
	})
}

func TestRoute_Matches(t *testing.T) {
	testCases := []struct {
		description string
		route       router.Route
		request     events.APIGatewayProxyRequest
		expected    bool
	}{
		{
			description: "Method and path matches",
			route: func() router.Route {
				route, _ := router.NewRoute(http.MethodPost, "\\/", nilHandler)
				return route
			}(),
			request:  events.APIGatewayProxyRequest{HTTPMethod: http.MethodPost, Path: "/"},
			expected: true,
		},
		{
			description: "Same path, different methods",
			route: func() router.Route {
				route, _ := router.NewRoute(http.MethodGet, "\\/", nilHandler)
				return route
			}(),
			request:  events.APIGatewayProxyRequest{HTTPMethod: http.MethodPost, Path: "/"},
			expected: false,
		},
		{
			description: "Same method, different paths",
			route: func() router.Route {
				route, _ := router.NewRoute(http.MethodPost, "\\/some-path", nilHandler)
				return route
			}(),
			request:  events.APIGatewayProxyRequest{HTTPMethod: http.MethodPost, Path: "/another-path"},
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			assert.Equal(
				t,
				testCase.expected,
				testCase.route.Matches(testCase.request),
			)
		})
	}
}
