package bridges

import (
	"context"
	"errors"
	routing "github.com/Napas/go-serverless-router"
	"github.com/Napas/go-serverless-router/routes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

const (
	awsRegionEuWest1 = "eu-west-1"
	sqsTargetArn     = "queue:arn"
	sqsQueueUrl      = "example.com/queue"
)

func Test_sqsBridge(t *testing.T) {
	nilLoggerMock := &loggerMock{}
	nilLoggerMock.On("Printf", mock.Anything, mock.Anything)

	nilRouter := &routerMock{}
	nilRouter.
		On("Handle", mock.Anything, mock.Anything).
		Return(nil, nil)

	t.Run("Passes SQS messages to the routing", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		bridgeCtx, cancel := context.WithCancel(ctx)

		sqsMock := &sqsMock{}
		sqsMock.
			On("ReceiveMessageWithContext",
				mock.Anything,
				&sqs.ReceiveMessageInput{
					QueueUrl: aws.String(sqsQueueUrl),
				},
				mock.Anything,
			).
			Once().
			Return(
				&sqs.ReceiveMessageOutput{
					Messages: []*sqs.Message{
						{
							MessageId:              aws.String("messageId"),
							ReceiptHandle:          aws.String("receiptHandle"),
							Body:                   aws.String("body"),
							MD5OfBody:              aws.String("md5OfBody"),
							MD5OfMessageAttributes: aws.String("md5OfMessageAttributes"),
							Attributes: map[string]*string{
								"attribute": aws.String("value"),
							},
							MessageAttributes: map[string]*sqs.MessageAttributeValue{
								"attribute": {
									StringValue: aws.String("value"),
								},
							},
						},
					},
				},
				nil,
			)
		sqsMock.
			On("ReceiveMessageWithContext",
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).
			Return(&sqs.ReceiveMessageOutput{}, nil)

		routerMock := &routerMock{}
		routerMock.
			On(
				"Handle",
				mock.MatchedBy(func(ctx aws.Context) bool {
					deadline, _ := ctx.Deadline()

					return assert.WithinDuration(t, time.Now().Add(time.Second*30), deadline, time.Second*2)
				}),
				map[string]interface{}{
					"Records": []interface{}{
						map[string]interface{}{
							"messageId":              aws.String("messageId"),
							"receiptHandle":          aws.String("receiptHandle"),
							"body":                   aws.String("body"),
							"md5OfBody":              aws.String("md5OfBody"),
							"md5OfMessageAttributes": aws.String("md5OfMessageAttributes"),
							"attributes": map[string]*string{
								"attribute": aws.String("value"),
							},
							"messageAttributes": map[string]*sqs.MessageAttributeValue{
								"attribute": {
									StringValue: aws.String("value"),
								},
							},
							"eventSourceARN": sqsTargetArn,
							"eventSource":    sqsQueueUrl,
							"awsRegion":      awsRegionEuWest1,
						},
					},
				}).
			Once().
			Return(nil, nil)

		bridge := NewSqsBridge(routerMock, sqsQueueUrl, sqsTargetArn, sqsMock, awsRegionEuWest1, nilLoggerMock)
		bridge.Run(bridgeCtx)

		time.Sleep(time.Millisecond * 100)

		sqsMock.AssertExpectations(t)
		routerMock.AssertExpectations(t)

		cancel()
	})

	t.Run("Log message if retrieving message failed", func(t *testing.T) {
		t.Parallel()

		logger := &loggerMock{}
		logger.
			On(
				"Printf",
				"Failed to consume message with error: %s",
				mock.MatchedBy(func(val []interface{}) bool {
					return val[0] == "error"
				}),
			).
			Once()

		logger.On("Printf", mock.Anything, mock.Anything)

		sqsMock := &sqsMock{}
		sqsMock.
			On("ReceiveMessageWithContext", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("error"))

		sqsMock.
			On("ReceiveMessageWithContext", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, nil)

		ctx := context.Background()
		bridgeCtx, cancel := context.WithCancel(ctx)

		bridge := NewSqsBridge(nilRouter, sqsQueueUrl, sqsTargetArn, sqsMock, awsRegionEuWest1, logger)
		bridge.Run(bridgeCtx)

		time.Sleep(time.Millisecond * 100)

		logger.AssertExpectations(t)

		cancel()
	})
}

type loggerMock struct {
	mock.Mock
	routing.Logger
}

func (m *loggerMock) Printf(val string, args ...interface{}) {
	m.Called(val, args)
}

type sqsMock struct {
	mock.Mock
	sqsiface.SQSAPI
}

func (m *sqsMock) ReceiveMessageWithContext(
	ctx aws.Context,
	input *sqs.ReceiveMessageInput,
	options ...request.Option,
) (output *sqs.ReceiveMessageOutput, err error) {
	args := m.Called(ctx, input, options)

	if args.Get(0) != nil {
		output = args.Get(0).(*sqs.ReceiveMessageOutput)
	}

	if args.Get(1) != nil {
		err = args.Error(1)
	}

	return output, err
}

type routerMock struct {
	mock.Mock
}

func (m *routerMock) AddRoute(route routes.Route) routing.Router {
	m.Called(route)

	return m
}

func (m *routerMock) Handle(ctx context.Context, event map[string]interface{}) (resp interface{}, err error) {
	args := m.Called(ctx, event)

	if args.Get(0) != nil {
		resp = args.Get(0)
	}

	if args.Get(1) != nil {
		resp = args.Error(1)
	}

	return resp, err
}
