package routes

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"regexp"
)

type SqsRoute struct {
	eventSourceArn *regexp.Regexp
	handler        SqsHandlerFunc
}

func NewSqsRoute(eventSourceArn string, handler SqsHandlerFunc) (*SqsRoute, error) {
	compiledEventSourceArn, err := regexp.Compile(eventSourceArn)

	if err != nil {
		return nil, RouteCompileError.Wrap(err, "Invalid regexp given")
	}

	return &SqsRoute{
		eventSourceArn: compiledEventSourceArn,
		handler:        handler,
	}, nil
}

func (route *SqsRoute) Matches(event map[string]interface{}) bool {
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

		eventSourceArn := recordVal["eventSourceARN"].(string)

		if !route.eventSourceArn.MatchString(eventSourceArn) {
			return false
		}

		return true
	}

	return true
}

func (route *SqsRoute) Handle(ctx context.Context, event map[string]interface{}) (interface{}, error) {
	jsonEvent, err := json.Marshal(event)

	if err != nil {
		return nil, RouteMarshalError.Wrap(err, "Failed to marshal event to JSON")
	}

	request := events.SQSEvent{}

	err = json.Unmarshal(jsonEvent, &request)

	if err != nil {
		return nil, RouteUnmarshalError.Wrap(err, "Failed to unmarshal request from the JSON")
	}

	return nil, route.handler(ctx, request)
}

func (*SqsRoute) HasResponse() bool {
	return false
}
