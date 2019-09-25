package bridges

import (
	"context"
	routing "github.com/Napas/go-serverless-router"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"time"
)

type sqsBridge struct {
	router    routing.Router
	queueUrl  string
	targetArn string
	sqs       sqsiface.SQSAPI
	awsRegion string
	logger    routing.Logger
}

func NewSqsBridge(
	r routing.Router,
	queueUrl string,
	targetArn string,
	sqs sqsiface.SQSAPI,
	awsRegion string,
	logger routing.Logger,
) Bridge {
	if logger == nil {
		logger = &routing.NilLogger{}
	}

	return &sqsBridge{
		router:    r,
		queueUrl:  queueUrl,
		targetArn: targetArn,
		sqs:       sqs,
		awsRegion: awsRegion,
		logger:    logger,
	}
}

// Run fetches messages from SQS and passes them to the routing.
// It's intended to be used for local development environments only.
func (bridge *sqsBridge) Run(ctx context.Context) {
	go func(ctx context.Context) {
		defer bridge.logger.Printf("Stopping SQS Bridge for: %s", bridge.queueUrl)

	consumer:
		for {
			select {
			case <-ctx.Done():
				break consumer
			default:
				err := bridge.receiveMessage(ctx)

				if err != nil {
					bridge.logger.Printf("Failed to consume message with error: %s", err.Error())
					// continue
				}
			}
		}
	}(ctx)
}

func (bridge *sqsBridge) receiveMessage(ctx context.Context) error {
	output, err := bridge.sqs.ReceiveMessageWithContext(
		ctx,
		&sqs.ReceiveMessageInput{
			QueueUrl: aws.String(bridge.queueUrl),
		},
	)

	if err != nil {
		return err
	}

	messagesCount := len(output.Messages)

	if messagesCount == 0 {
		return nil
	}

	bridge.logger.Printf(
		"Received %d messages from the %s queue, passing them to the routing",
		messagesCount,
		bridge.queueUrl,
	)

	records := make([]interface{}, messagesCount)

	for i, message := range output.Messages {
		records[i] = map[string]interface{}{
			"messageId":              message.MessageId,
			"receiptHandle":          message.ReceiptHandle,
			"body":                   message.Body,
			"md5OfBody":              message.MD5OfBody,
			"md5OfMessageAttributes": message.MD5OfMessageAttributes,
			"attributes":             message.Attributes,
			"messageAttributes":      message.MessageAttributes,
			"eventSourceARN":         bridge.targetArn,
			"eventSource":            bridge.queueUrl,
			"awsRegion":              bridge.awsRegion,
		}
	}

	event := map[string]interface{}{
		"Records": records,
	}

	reqCtx, _ := context.WithDeadline(ctx, time.Now().Add(time.Second*30))
	_, err = bridge.router.Handle(reqCtx, event)

	return err
}
