package sqs

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Client interface {
	SendMessage(ctx context.Context, message string) error
}

type ClientImpl struct {
	client   *sqs.Client
	queueUrl string
}

func New(client *sqs.Client, queueUrl string) *ClientImpl {
	log.Printf("Creating SQS client with Queue URL: %s", queueUrl)
	return &ClientImpl{
		client:   client,
		queueUrl: queueUrl,
	}
}

func (c *ClientImpl) SendMessage(ctx context.Context, message string) error {
	log.Printf("Sending message to SQS: %s", message)
	log.Printf("Queue URL: %s", c.queueUrl)
	_, err := c.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &c.queueUrl,
		MessageBody: &message,
	})

	return err
}
