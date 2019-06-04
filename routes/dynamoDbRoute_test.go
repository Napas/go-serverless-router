package routes_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/joomcode/errorx"

	"github.com/aws/aws-lambda-go/events"

	"github.com/Napas/go-serverless-router/routes"
	"github.com/stretchr/testify/assert"
)

func Test_DynamoDbRoute(t *testing.T) {
	voidHandler := func(ctx context.Context, request events.DynamoDBEvent) {}

	t.Run("NewDynamoDbRoute", func(t *testing.T) {
		t.Run("Returns an error if invalid regexp is passed", func(t *testing.T) {
			_, err := routes.NewDynamoDbRoute("[invalid regexp", voidHandler)

			assert.Error(t, err)
			assert.True(t, err.(*errorx.Error).IsOfType(routes.RouteCompileError))
		})
	})

	t.Run("HasResponse returns false", func(t *testing.T) {
		router, err := routes.NewDynamoDbRoute(".*", voidHandler)

		assert.Nil(t, err)
		assert.False(t, router.HasResponse())
	})

	t.Run("Matches", func(t *testing.T) {
		t.Run("Returns false", func(t *testing.T) {
			t.Run("If Records key is not set", func(t *testing.T) {
				router, err := routes.NewDynamoDbRoute(".*", voidHandler)

				assert.Nil(t, err)
				assert.False(t, router.Matches(map[string]interface{}{}))
			})

			t.Run("If Records is not a slice of []map[string]interface{}", func(t *testing.T) {
				router, err := routes.NewDynamoDbRoute(".*", voidHandler)

				assert.Nil(t, err)
				assert.False(t, router.Matches(map[string]interface{}{
					"Records": []string{},
				}))
			})

			t.Run("If Records is an empty slice", func(t *testing.T) {
				router, err := routes.NewDynamoDbRoute(".*", voidHandler)

				assert.Nil(t, err)
				assert.False(t, router.Matches(map[string]interface{}{
					"Records": []map[string]interface{}{},
				}))
			})

			t.Run("If dynamodb in the record is empty", func(t *testing.T) {
				router, err := routes.NewDynamoDbRoute("/something/", voidHandler)

				assert.Nil(t, err)
				assert.False(t, router.Matches(map[string]interface{}{
					"Records": []map[string]interface{}{
						{
							"eventSourceARN": "something",
						},
					},
				}))
			})

			t.Run("If eventSourceARN doe not match in the first record", func(t *testing.T) {
				router, err := routes.NewDynamoDbRoute("another", voidHandler)

				assert.Nil(t, err)
				assert.False(t, router.Matches(map[string]interface{}{
					"Records": []map[string]interface{}{
						{
							"eventSourceARN": "something",
							"dynamodb": map[string]interface{}{
								"key": "value",
							},
						},
					},
				}))
			})
		})

		t.Run("Should be true", func(t *testing.T) {
			router, err := routes.NewDynamoDbRoute("^arn:aws:dynamodb:us-east-1:.*:table\\/.*\\/stream.*$", voidHandler)
			request := map[string]interface{}{}
			requestJson := `{
  "Records": [
    {
      "dynamodb": {
        "ApproximateCreationDateTime": 1559762747,
        "Keys": {
          "id": {
            "S": "d21eae87-47b1-46b4-b74b-30b63f3b971c"
          }
        },
        "SequenceNumber": "88625000000000079984566295",
        "SizeBytes": 38,
        "StreamViewType": "KEYS_ONLY"
      },
      "eventSourceARN": "arn:aws:dynamodb:us-east-1:1111111:table/table/stream/2019-05-31T17:47:06.850"
    }
  ]
}`
			assert.Nil(t, err)

			err = json.Unmarshal([]byte(requestJson), &request)

			assert.Nil(t, err)
			assert.True(t, router.Matches(request))
		})
	})

	t.Run("Handle", func(t *testing.T) {
		t.Run("Passes correct data to the handler", func(t *testing.T) {
			requestContext := context.TODO()

			route, err := routes.NewDynamoDbRoute(
				".*",
				func(ctx context.Context, request events.DynamoDBEvent) {
					assert.Equal(t, requestContext, ctx)
					assert.Equal(t, "eventID", request.Records[0].EventID)
				},
			)

			assert.Nil(t, err)

			event := make(map[string]interface{})
			event["Records"] = []map[string]interface{}{
				{
					"eventID": "eventID",
				},
			}

			route.Handle(requestContext, event)
		})
	})
}
