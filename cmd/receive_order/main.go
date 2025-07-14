package main

import (
	"context"
	"encoding/json"
	"os"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Order struct {
	OrderID      string `json:"orderId"`
	CustomerName string `json:"customerName"`
}

var (
	sess      = session.Must(session.NewSession())
	sqsClient = sqs.New(sess)
	queueURL  = os.Getenv("QUEUE_URL")
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var order Order
	err := json.Unmarshal([]byte(req.Body), &order)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid input"}, nil
	}
	body, _ := json.Marshal(order)
	_, err = sqsClient.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(body)),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Failed to send message"}, nil
	}
	return events.APIGatewayProxyResponse{StatusCode: 202, Body: "Order received"}, nil
}

func main() {
	lambda.Start(handler)
}
