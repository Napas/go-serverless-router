package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
)

type ApiGatewayRoute struct {
	path       *regexp.Regexp
	httpMethod string
	handler    ApiGatewayHandlerFunc
}

func NewApiGatewayRoute(
	path string,
	httpMethod string,
	handler ApiGatewayHandlerFunc,
) (*ApiGatewayRoute, error) {
	compiledPath, err := regexp.Compile(path)

	if err != nil {
		return nil, RouteCompileError.Wrap(err, "Invalid regexp given")
	}

	return &ApiGatewayRoute{
		path:       compiledPath,
		httpMethod: httpMethod,
		handler:    handler,
	}, nil
}

func (route *ApiGatewayRoute) Matches(event map[string]interface{}) bool {
	if event["httpMethod"] != route.httpMethod {
		return false
	}

	if !route.path.MatchString(event["path"].(string)) {
		return false
	}

	return true
}

func (route *ApiGatewayRoute) Handle(ctx context.Context, event map[string]interface{}) (interface{}, error) {
	jsonEvent, err := json.Marshal(event)

	if err != nil {
		return events.APIGatewayProxyResponse{}, RouteMarshalError.Wrap(err, "Failed to marshal event to JSON")
	}

	request := events.APIGatewayProxyRequest{}
	err = json.Unmarshal(jsonEvent, &request)

	if err != nil {
		return events.APIGatewayProxyResponse{}, RouteUnmarshalError.Wrap(err, "Failed to unmarshal event from JSON")
	}

	return route.handler(ctx, request)
}

func (*ApiGatewayRoute) HasResponse() bool {
	return true
}

func (route *ApiGatewayRoute) String() string {
	return fmt.Sprintf("API Gateway route: %s %s", route.httpMethod, route.path.String())
}
