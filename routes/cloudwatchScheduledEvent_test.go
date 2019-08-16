package routes_test

import (
	"context"
	"github.com/Napas/go-serverless-router/routes"
	"github.com/aws/aws-lambda-go/events"
	"github.com/joomcode/errorx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestScheduledEventRoute(t *testing.T) {
	t.Parallel()

	nilHandler := func(ctx context.Context, request events.CloudWatchEvent) error {
		return nil
	}

	t.Run("NewCloudwatchScheduledEventRoute", func(t *testing.T) {
		t.Run("Returns an error if invalid regexp is passed", func(t *testing.T) {
			_, err := routes.NewCloudwatchScheduledEventRoute([]string{"[invalid regexp"}, nilHandler)

			assert.Error(t, err)
			assert.True(t, err.(*errorx.Error).IsOfType(routes.RouteCompileError))
		})
	})

	t.Run("HasResponse returns false", func(t *testing.T) {
		route, err := routes.NewCloudwatchScheduledEventRoute([]string{".*"}, nilHandler)

		require.Nil(t, err)
		assert.False(t, route.HasResponse())
	})

	t.Run("Matches", func(t *testing.T) {
		t.Run("Returns false", func(t *testing.T) {
			t.Run("If detail-type is equal to the 'Scheduled event'", func(t *testing.T) {
				route, err := routes.NewCloudwatchScheduledEventRoute([]string{".*"}, nilHandler)

				require.Nil(t, err)
				assert.False(t, route.Matches(map[string]interface{}{}))
			})

			t.Run("If resources is not set", func(t *testing.T) {
				route, err := routes.NewCloudwatchScheduledEventRoute([]string{".*"}, nilHandler)

				require.Nil(t, err)
				assert.False(t, route.Matches(map[string]interface{}{
					"detail-type": "Scheduled Event",
				}))
			})

			t.Run("If not all required resource ARNs matches event ARNs", func(t *testing.T) {
				testCases := []struct {
					description  string
					event        map[string]interface{}
					resourceArns []string
				}{
					{
						description: "Not an equal number of resources",
						event: map[string]interface{}{
							"detail-type": "Scheduled Event",
							"resources":   []string{"arn:1"},
						},
						resourceArns: []string{"arn:1", "arn:2"},
					},
					{
						description: "Has different ARNs",
						event: map[string]interface{}{
							"detail-type": "Scheduled Event",
							"resources":   []string{"arn:1"},
						},
						resourceArns: []string{"arn:2"},
					},
				}

				for _, testCase := range testCases {
					t.Run(testCase.description, func(t *testing.T) {
						route, err := routes.NewCloudwatchScheduledEventRoute(
							testCase.resourceArns,
							nilHandler,
						)

						require.Nil(t, err)
						assert.False(t, route.Matches(testCase.event))
					})
				}
			})
		})

		t.Run("Returns true otherwise", func(t *testing.T) {
			route, err := routes.NewCloudwatchScheduledEventRoute([]string{".*"}, nilHandler)

			require.Nil(t, err)
			assert.True(t, route.Matches(map[string]interface{}{
				"detail-type": "Scheduled Event",
				"resources":   []string{"arn:1"},
			}))
		})
	})

	t.Run("Handle", func(t *testing.T) {
		t.Run("Passes correct data to the handler", func(t *testing.T) {
			requestCtx := context.TODO()

			route, err := routes.NewCloudwatchScheduledEventRoute(
				[]string{".*"},
				func(ctx context.Context, request events.CloudWatchEvent) error {
					assert.Equal(t, requestCtx, ctx)
					assert.Equal(t, "ID", request.ID)

					return nil
				},
			)

			require.Nil(t, err)

			resp, err := route.Handle(
				requestCtx,
				map[string]interface{}{
					"detail-type": "Scheduled Event",
					"resources":   []string{"arn:1"},
					"ID":          "ID",
				},
			)

			assert.Nil(t, err)
			assert.Nil(t, resp)
		})
	})
}
