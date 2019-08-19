package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"regexp"
	"strings"
)

const (
	scheduledEventName = "Scheduled Event"
)

type CloudwatchScheduledEventRoute struct {
	resourceArns resourceArnsRegexps
	handler      CloudWatchScheduledEventHandlerFunc
}

func NewCloudwatchScheduledEventRoute(
	resourceArns []string,
	handler CloudWatchScheduledEventHandlerFunc,
) (*CloudwatchScheduledEventRoute, error) {
	route := &CloudwatchScheduledEventRoute{
		handler:      handler,
		resourceArns: resourceArnsRegexps{},
	}

	for _, resourceArn := range resourceArns {
		compiledResourceArn, err := regexp.Compile(resourceArn)

		if err != nil {
			return nil, RouteCompileError.Wrap(err, "Invalid regexp given")
		}

		route.resourceArns = append(route.resourceArns, compiledResourceArn)
	}

	return route, nil
}

func (route *CloudwatchScheduledEventRoute) Matches(event map[string]interface{}) bool {
	if event["detail-type"] != scheduledEventName {
		return false
	}

	resourceArns, ok := event["resources"].([]interface{})

	if !ok {
		return false
	}

	if len(resourceArns) != len(route.resourceArns) {
		return false
	}

	for _, resourceArn := range resourceArns {
		resourceArnStr, ok := resourceArn.(string)

		if !ok || !route.resourceArns.has(resourceArnStr) {
			return false
		}
	}

	return true
}

func (route *CloudwatchScheduledEventRoute) Handle(ctx context.Context, event map[string]interface{}) (interface{}, error) {
	jsonEvent, err := json.Marshal(event)

	if err != nil {
		return nil, RouteMarshalError.Wrap(err, "Failed to marshal event to JSON")
	}

	request := events.CloudWatchEvent{}

	err = json.Unmarshal(jsonEvent, &request)

	if err != nil {
		return nil, RouteUnmarshalError.Wrap(err, "Failed to unmarshal request from the JSON")
	}

	err = route.handler(ctx, request)

	return err, err
}

func (*CloudwatchScheduledEventRoute) HasResponse() bool {
	return false
}

func (route *CloudwatchScheduledEventRoute) String() string {
	resources := []string{}

	for _, resource := range route.resourceArns {
		resources = append(resources, resource.String())
	}

	return fmt.Sprintf("CloudWatch scheduled event for resources: %s", strings.Join(resources, ", "))
}

type resourceArnsRegexps []*regexp.Regexp

func (arns resourceArnsRegexps) has(arn string) bool {
	for _, expectedArn := range arns {
		if expectedArn.MatchString(arn) {
			return true
		}
	}

	return false
}
