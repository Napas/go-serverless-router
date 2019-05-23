package goserverlessrouter_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/Napas/go-serverless-router"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

type handlerMock struct {
	mock.Mock
}

func (handler *handlerMock) Handle(
	ctx context.Context,
	request events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	args := handler.Called(ctx, request)

	return args.Get(0).(events.APIGatewayProxyResponse), args.Error(1)
}

func TestRouter_Handle(t *testing.T) {
	t.Run("Return response with 404 if no routes matches", func(t *testing.T) {
		router := goserverlessrouter.NewRouter()
		response, _ := router.Handle(context.TODO(), events.APIGatewayProxyRequest{})

		assert.Equal(t, http.StatusNotFound, response.StatusCode)
	})

	t.Run("Returns response from the matched route handler", func(t *testing.T) {
		ctx := context.TODO()
		request := events.APIGatewayProxyRequest{
			HTTPMethod: http.MethodGet,
			Path:       "/",
		}
		expectedResponse := events.APIGatewayProxyResponse{
			StatusCode: http.StatusNoContent,
		}

		handlerMock := new(handlerMock)
		handlerMock.On("Handle", ctx, request).Return(expectedResponse, nil)

		route, _ := goserverlessrouter.NewRoute(http.MethodGet, "/", handlerMock.Handle)
		router := goserverlessrouter.NewRouter()
		router.AddRoute(route)
		response, err := router.Handle(ctx, request)

		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, response)
	})

	t.Run("Returns error from the matched route handler", func(t *testing.T) {
		ctx := context.TODO()
		request := events.APIGatewayProxyRequest{
			HTTPMethod: http.MethodGet,
			Path:       "/",
		}
		expectedErr := errors.New("expected error from the handler")

		handlerMock := new(handlerMock)
		handlerMock.On("Handle", ctx, request).Return(events.APIGatewayProxyResponse{}, expectedErr)

		route, _ := goserverlessrouter.NewRoute(http.MethodGet, "/", handlerMock.Handle)
		router := goserverlessrouter.NewRouter()
		router.AddRoute(route)
		_, err := router.Handle(ctx, request)

		assert.Equal(t, expectedErr, err)
	})
}
