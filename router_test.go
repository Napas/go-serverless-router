package routing_test

import (
	"context"
	"errors"
	"testing"

	"github.com/joomcode/errorx"

	goserverlessrouter "github.com/Napas/go-serverless-router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type routeMock struct {
	mock.Mock
}

func (route *routeMock) Matches(event map[string]interface{}) bool {
	args := route.Called(event)

	return args.Bool(0)
}

func (route *routeMock) Handle(ctx context.Context, event map[string]interface{}) (interface{}, error) {
	args := route.Called(ctx, event)

	return args.Get(0), args.Error(1)
}

func (route *routeMock) HasResponse() bool {
	args := route.Called()

	return args.Bool(0)
}

func Test_Router(t *testing.T) {
	t.Parallel()

	t.Run("Returning an error if route is not found", func(t *testing.T) {
		router := goserverlessrouter.New()
		err1, err2 := router.Handle(context.TODO(), map[string]interface{}{})

		assert.Error(t, err1.(error))
		assert.Error(t, err2)

		assert.True(t, err1.(*errorx.Error).IsOfType(goserverlessrouter.RouterRouteNotFoundError))
		assert.True(t, err2.(*errorx.Error).IsOfType(goserverlessrouter.RouterRouteNotFoundError))
	})

	t.Run("Calls Handle function on matching route", func(t *testing.T) {
		ctx := context.TODO()
		event := map[string]interface{}{}

		route1 := &routeMock{}
		route1.
			On("Matches", event).
			Once().
			Return(false)

		route2 := &routeMock{}
		route2.
			On("Matches", event).
			Once().
			Return(true)

		route2.
			On("Handle", ctx, event).
			Once().
			Return(nil, nil)

		route2.
			On("HasResponse").
			Once().
			Return(true)

		router := goserverlessrouter.New()
		router.
			AddRoute(route1).
			AddRoute(route2)

		router.Handle(ctx, event)

		route1.AssertExpectations(t)
		route2.AssertExpectations(t)
	})

	t.Run("Returns response and error if HasResponse returns true", func(t *testing.T) {
		expectedResponse := make(map[string]interface{})
		expectedResponse["Body"] = "response"

		expectedError := errors.New("error")

		route := &routeMock{}
		route.
			On("Matches", mock.Anything).
			Return(true)

		route.
			On("Handle", mock.Anything, mock.Anything).
			Return(expectedResponse, expectedError)

		route.
			On("HasResponse").
			Return(true)

		router := goserverlessrouter.New()
		router.AddRoute(route)

		resp, err := router.Handle(context.TODO(), map[string]interface{}{})

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, expectedResponse, resp)
	})

	t.Run("Returns error af both return arguments if HasResponse is false", func(t *testing.T) {
		response := make(map[string]interface{})
		response["Body"] = "response"

		expectedError := errors.New("error")

		route := &routeMock{}
		route.
			On("Matches", mock.Anything).
			Return(true)

		route.
			On("Handle", mock.Anything, mock.Anything).
			Return(response, expectedError)

		route.
			On("HasResponse").
			Return(false)

		router := goserverlessrouter.New()
		router.AddRoute(route)

		resp, err := router.Handle(context.TODO(), map[string]interface{}{})

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Error(t, resp.(error))
		assert.Equal(t, expectedError, resp.(error))
	})
}
