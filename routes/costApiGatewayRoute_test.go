package routes_test

import (
	"context"
	"github.com/Napas/go-serverless-router/routes"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CorsApiGatewayRoute(t *testing.T) {
	t.Parallel()

	t.Run("NewCorsApiGatewayRoute", func(t *testing.T) {
		t.Run("Returns an error if path does not compile to regexp", func(t *testing.T) {
			_, err := routes.NewCorsApiGatewayRoute(
				"[invalid regexp",
				"",
				nil,
				nil,
			)

			assert.Error(t, err)
		})
	})

	t.Run("handler", func(t *testing.T) {
		t.Run("Adds Access-Control-Allow-Origin header", func(t *testing.T) {
			handler, err := routes.NewCorsApiGatewayRoute("/*./", "http://example.com", nil, nil)

			assert.Nil(t, err)

			resp, err := handler.Handle(context.TODO(), map[string]interface{}{"path": "/"})

			assert.Nil(t, err)
			assert.IsType(t, resp, events.APIGatewayProxyResponse{})

			apiResp := resp.(events.APIGatewayProxyResponse)

			assert.Equal(t, "http://example.com", apiResp.Headers["Access-Control-Allow-Origin"])
		})

		t.Run("Adds Access-Control-Allow-Methods header", func(t *testing.T) {
			handler, err := routes.NewCorsApiGatewayRoute("/*./", "", []string{"GET", "OPTIONS"}, nil)

			assert.Nil(t, err)

			resp, err := handler.Handle(context.TODO(), map[string]interface{}{"path": "/"})

			assert.Nil(t, err)
			assert.IsType(t, resp, events.APIGatewayProxyResponse{})

			apiResp := resp.(events.APIGatewayProxyResponse)

			assert.Equal(t, "GET, OPTIONS", apiResp.Headers["Access-Control-Allow-Methods"])
		})

		t.Run("Adds OPTIONS if it does not exist into Access-Control-Allow-Methods", func(t *testing.T) {
			handler, err := routes.NewCorsApiGatewayRoute("/*./", "", []string{"GET"}, nil)

			assert.Nil(t, err)

			resp, err := handler.Handle(context.TODO(), map[string]interface{}{"path": "/"})

			assert.Nil(t, err)
			assert.IsType(t, resp, events.APIGatewayProxyResponse{})

			apiResp := resp.(events.APIGatewayProxyResponse)

			assert.Equal(t, "GET, OPTIONS", apiResp.Headers["Access-Control-Allow-Methods"])
		})

		t.Run("Add Access-Control-Allow-Headers header", func(t *testing.T) {
			handler, err := routes.NewCorsApiGatewayRoute("/*./", "", nil, []string{"Accept", "Content-Type"})

			assert.Nil(t, err)

			resp, err := handler.Handle(context.TODO(), map[string]interface{}{"path": "/"})

			assert.Nil(t, err)
			assert.IsType(t, resp, events.APIGatewayProxyResponse{})

			apiResp := resp.(events.APIGatewayProxyResponse)

			assert.Equal(t, "Accept, Content-Type", apiResp.Headers["Access-Control-Allow-Headers"])
		})
	})
}
