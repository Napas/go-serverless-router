package bridges

import (
	"context"
	router "github.com/Napas/go-serverless-router"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"time"
)

type sqsBridge struct {
	router    router.Router
	queueUrl  string
	targetArn string
	sqs       sqsiface.SQSAPI
	awsRegion string
	logger    router.Logger
}

func NewSqsBridge(
	r router.Router,
	queueUrl string,
	targetArn string,
	sqs sqsiface.SQSAPI,
	awsRegion string,
	logger router.Logger,
) Bridge {
	if logger == nil {
		logger = &router.NilLogger{}
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

// Run fetches messages from SQS and passes them to the router.
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
		"Received %d messages from the %s queue, passing them to the router",
		messagesCount,
		bridge.queueUrl,
	)

	records := make([]map[string]interface{}, messagesCount)

	for i, message := range output.Messages {
		records[i] = map[string]interface{}{
			"MessageId":              message.MessageId,
			"ReceiptHandle":          message.ReceiptHandle,
			"Body":                   message.Body,
			"Md5OfBody":              message.MD5OfBody,
			"Md5OfMessageAttributes": message.MD5OfMessageAttributes,
			"Attributes":             message.Attributes,
			"MessageAttributes":      message.MessageAttributes,
			"EventSourceARN":         bridge.targetArn,
			"EventSource":            bridge.queueUrl,
			"AWSRegion":              bridge.awsRegion,
		}
	}

	event := map[string]interface{}{
		"Records": records,
	}

	reqCtx, _ := context.WithDeadline(ctx, time.Now().Add(time.Second*30))
	_, err = bridge.router.Handle(reqCtx, event)

	return err
}
