package routes_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/Napas/go-serverless-router/routes"

	"github.com/aws/aws-lambda-go/events"

	"github.com/stretchr/testify/assert"
)

func Test_ApiGatewayRoute(t *testing.T) {
	voidHandler := func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return events.APIGatewayProxyResponse{}, nil
	}

	t.Run("NewApiGatewayRoute", func(t *testing.T) {
		t.Run("Returns an error if path does not compile to regexp", func(t *testing.T) {
			_, err := routes.NewApiGatewayRoute(
				"[invalid regexp",
				http.MethodGet,
				voidHandler,
			)

			assert.Error(t, err)
		})
	})

	t.Run("Matches", func(t *testing.T) {
		t.Run("Returns false", func(t *testing.T) {
			t.Run("If http methods does not match", func(t *testing.T) {
				testCases := []struct {
					name            string
					httpMethod      string
					eventHttpMethod string
				}{
					{
						"Empty http method",
						"",
						http.MethodGet,
					},
					{
						"non empty http method",
						http.MethodPut,
						http.MethodGet,
					},
				}

				for _, testCase := range testCases {
					t.Run(testCase.name, func(t *testing.T) {
						route, err := routes.NewApiGatewayRoute("/*./", testCase.httpMethod, voidHandler)

						assert.NoError(t, err)

						event := make(map[string]interface{})
						event["httpMethod"] = testCase.eventHttpMethod

						assert.False(t, route.Matches(event))
					})
				}
			})

			t.Run("If path does not match regexp", func(t *testing.T) {
				route, err := routes.NewApiGatewayRoute("/abc", http.MethodGet, voidHandler)

				assert.NoError(t, err)

				event := make(map[string]interface{})
				event["httpMethod"] = http.MethodGet
				event["path"] = "/non-matching-path"

				assert.False(t, route.Matches(event))
			})
		})

		t.Run("Returns true if http method and path matches", func(t *testing.T) {
			route, err := routes.NewApiGatewayRoute("/path", http.MethodPost, voidHandler)

			assert.NoError(t, err)

			event := make(map[string]interface{})
			event["httpMethod"] = http.MethodPost
			event["path"] = "\\/path"

			assert.True(t, route.Matches(event))
		})
	})

	t.Run("Handle", func(t *testing.T) {
		t.Run("Passes correct data to the handler", func(t *testing.T) {
			requestContext := context.TODO()

			route, err := routes.NewApiGatewayRoute(
				"\\/path",
				http.MethodPost,
				func(ctx context.Context, request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, e error) {
					assert.Equal(t, requestContext, ctx)
					assert.Equal(t, http.MethodPost, request.HTTPMethod)
					assert.Equal(t, "/path", request.Path)

					return events.APIGatewayProxyResponse{}, nil
				},
			)

			assert.NoError(t, err)

			event := make(map[string]interface{})
			event["httpMethod"] = http.MethodPost
			event["path"] = "/path"

			route.Handle(requestContext, event)
		})

		t.Run("Returns data from the handler", func(t *testing.T) {
			response := events.APIGatewayProxyResponse{
				Body: "response",
			}
			responseErr := errors.New("Response error")

			route, err := routes.NewApiGatewayRoute(
				"\\/path",
				http.MethodPost,
				func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
					return response, responseErr
				},
			)

			assert.NoError(t, err)

			event := make(map[string]interface{})
			event["httpMethod"] = http.MethodPost
			event["path"] = "/path"

			resp, err := route.Handle(context.TODO(), event)

			assert.Equal(t, response, resp)
			assert.Equal(t, responseErr, err)
		})
	})
}
