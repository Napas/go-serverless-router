package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
)

type DynamoDbRoute struct {
	eventSourceArn *regexp.Regexp
	handler        DynamoDbHandlerFunc
}

func NewDynamoDbRoute(
	eventSourceArn string,
	handler DynamoDbHandlerFunc,
) (*DynamoDbRoute, error) {
	compiledEventSourceArn, err := regexp.Compile(eventSourceArn)

	if err != nil {
		return nil, RouteCompileError.Wrap(err, "Invalid regexp given")
	}

	return &DynamoDbRoute{
		eventSourceArn: compiledEventSourceArn,
		handler:        handler,
	}, nil
}

func (route *DynamoDbRoute) Matches(event map[string]interface{}) bool {
	if event["Records"] == nil {
		return false
	}

	records, ok := event["Records"].([]interface{})

	if !ok {
		return false
	}

	if len(records) == 0 {
		return false
	}

	for _, record := range records {
		recordVal, ok := record.(map[string]interface{})

		if !ok {
			return false
		}

		dynamodb, ok := recordVal["dynamodb"].(map[string]interface{})

		if !ok || dynamodb == nil || recordVal["dynamodb"] == nil {
			return false
		}

		eventSourceArn := recordVal["eventSourceARN"].(string)

		if !route.eventSourceArn.MatchString(eventSourceArn) {
			return false
		}

		return true
	}

	return true
}

func (route *DynamoDbRoute) Handle(ctx context.Context, event map[string]interface{}) (interface{}, error) {
	jsonEvent, err := json.Marshal(event)

	if err != nil {
		return nil, RouteMarshalError.Wrap(err, "Failed to marshal event to JSON")
	}

	request := events.DynamoDBEvent{}

	err = json.Unmarshal(jsonEvent, &request)

	if err != nil {
		return nil, RouteUnmarshalError.Wrap(err, "Failed to unmarshal request from the JSON")
	}

	route.handler(ctx, request)

	return nil, nil
}

func (*DynamoDbRoute) HasResponse() bool {
	return false
}

func (route *DynamoDbRoute) String() string {
	return fmt.Sprintf("Dynamo db event for %s", route.eventSourceArn.String())
}
