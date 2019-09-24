package routes_test

import (
	"context"
	"encoding/json"
	"github.com/Napas/go-serverless-router/routes"
	"github.com/aws/aws-lambda-go/events"
	"github.com/joomcode/errorx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SqsRoute(t *testing.T) {
	t.Parallel()

	nilHandler := func(_ context.Context, _ events.SQSEvent) error {
		return nil
	}

	t.Run("NewSqsRoute", func(t *testing.T) {
		t.Run("Returns error on the invalid eventSourceArn regexp", func(t *testing.T) {
			_, err := routes.NewSqsRoute("[[", nilHandler)

			assert.Error(t, err)
			assert.True(t, err.(*errorx.Error).IsOfType(routes.RouteCompileError))
		})
	})

	t.Run("HasResponse", func(t *testing.T) {
		t.Run("Returns false", func(t *testing.T) {
			sqsRoute, _ := routes.NewSqsRoute("/.*/", nilHandler)

			assert.False(t, sqsRoute.HasResponse())
		})
	})

	t.Run("Matches", func(t *testing.T) {
		testCases := []struct {
			request  string
			expected bool
		}{
			{
				request:  "",
				expected: false,
			},
			{
				request:  `{"Records":[]}`,
				expected: false,
			},
			{
				request:  `{"Records":[{"eventSourceARN": ""}]}`,
				expected: false,
			},
			{
				request:  `{"Records":[{"eventSourceARN": "non-matching"}]}`,
				expected: false,
			},
			{
				request:  `{"Records":[{"eventSourceARN": "arn:aws:sqs:us-east-2:123456789012:my-queue"}]}`,
				expected: true,
			},
		}

		for i, testCase := range testCases {
			t.Run(string(i), func(t *testing.T) {
				sqsRoute, _ := routes.NewSqsRoute("^arn:aws:sqs:us-east-2:123456789012:my-queue$", nilHandler)
				req := map[string]interface{}{}
				json.Unmarshal([]byte(testCase.request), &req)

				assert.Equal(t, testCase.expected, sqsRoute.Matches(req))
			})
		}
	})

	t.Run("Handle", func(t *testing.T) {
		t.Run("Passes correct data to the handler", func(t *testing.T) {
			requestContext := context.TODO()

			route, err := routes.NewSqsRoute(
				".*",
				func(ctx context.Context, request events.SQSEvent) error {
					assert.Equal(t, requestContext, ctx)
					assert.Equal(t, "eventSourceARN", request.Records[0].EventSourceARN)

					return nil
				},
			)

			assert.Nil(t, err)

			event := make(map[string]interface{})
			event["Records"] = []map[string]interface{}{
				{
					"eventSourceARN": "eventSourceARN",
				},
			}

			route.Handle(requestContext, event)
		})
	})
}
